package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"

	"stockbit-challenge/application"
	"stockbit-challenge/controller"
)

const (
	clientMode   = "client"
	consumerMode = "consumer"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		args = os.Args[1:]
		mode string
	)

	if len(args) > 0 {
		mode = os.Args[1]
	}

	ctx, cancel := context.WithCancel(context.Background())
	app := application.SetupApplication(ctx)
	defer app.Close()

	// Setup dependency injection
	dep := application.SetupDependency(app)

	// Channel for OS interruption
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// Goroutine to cancel when Os interrupt happens
	go func() {
		fmt.Println(fmt.Sprintf("system call: %+v", <-c))
		cancel()
	}()

	switch mode {
	case clientMode:
		httpServer := serveHTTP(app, dep)
		defer func(httpServer *http.Server) {
			_ = httpServer.Close()
		}(httpServer)

		<-app.Context.Done()
	case consumerMode:
		consumeKafkaMessages(app, dep)
	default:
		httpServer := serveHTTP(app, dep)
		defer func(httpServer *http.Server) {
			_ = httpServer.Close()
		}(httpServer)

		consumeKafkaMessages(app, dep)

		<-app.Context.Done()
	}
}

func serveHTTP(app *application.App, dep *application.Dependency) *http.Server {
	r := chi.NewRouter()

	trxController := &controller.TransactionController{
		TransactionFeedService: dep.TransactionFeedService,
	}

	r.Post("/publish/transaction", trxController.UploadTransactions)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%v", app.Config.Const.HTTPPort),
		Handler:        r,
		ReadTimeout:    time.Duration(app.Config.Const.ShortTimeout) * time.Second,
		WriteTimeout:   time.Duration(app.Config.Const.ShortTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func(hs *http.Server) {
		err := hs.ListenAndServe()
		if err != nil {
			fmt.Println("Failed to serve http server")
		}
	}(s)

	return s
}

func consumeKafkaMessages(app *application.App, dep *application.Dependency) {
	topics := app.Config.Kafka.ConsumerTopics
	consumerHandler := &controller.ConsumerHandler{
		TransactionFeed: dep.TransactionFeedService,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(t string, h func([]byte) (bool, error)) {
		fmt.Println(fmt.Sprintf("creating consumer for topic %s", t))
		kafkaListener(app, t, h)
		wg.Done()
	}(topics["transaction"], consumerHandler.Transaction)
}

func kafkaListener(app *application.App, topic string, messageHandler func([]byte) (bool, error)) {
	msgInfo := fmt.Sprintf("Kafka Listener(%s) START", topic)
	fmt.Println(msgInfo)

loop:
	for {
		select {
		//this block will close current consumer, triggered by context done
		case <-app.Context.Done():
			fmt.Println(fmt.Sprintf("consumer stop signal %s", topic))
			break loop
		default:
			msg, err := app.Consumer.Consume(topic)
			if err != nil {
				if err == io.EOF || err == context.Canceled {
					fmt.Println(fmt.Sprintf("shutting down kafka consumers %s", topic))
					break loop
				}

				fmt.Println(fmt.Sprintf("error when consuming kafka message %s", topic))
				continue
			} else if msg.Value == nil {
				// no message, no error. skip
				continue
			}
			// Process the message.
			fmt.Println(fmt.Sprintf("kafka message consumed %s[%d]%d", topic, msg.Partition, msg.Offset))

			if _, err := messageHandler(msg.Value.([]byte)); err != nil { //error on adapter or post-operation
				fmt.Println(fmt.Sprintf("error processing message %s[%d]%d", topic, msg.Partition, msg.Offset))
			}
		}
	}

	_ = app.Consumer.Close()
	fmt.Println(fmt.Sprintf("Listener (%s) STOP", topic))
}
