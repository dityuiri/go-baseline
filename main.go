package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"

	"stockbit-challenge/application"
	"stockbit-challenge/controller"
	pb "stockbit-challenge/proto/proto-golang/stock"
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
	app := application.SetupApplication()

	// Setup dependency injection
	dep := application.SetupDependency(app)

	// Channel for OS interruption
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	// Goroutine to cancel when Os interrupt happens
	go func() {
		log.Printf("system call: %+v", <-c)
		cancel()
	}()

	switch mode {
	case clientMode:
		httpServer := serveHTTP(app, dep)
		defer func(httpServer *http.Server) {
			_ = httpServer.Close()
		}(httpServer)

		grpcServer := serveGRPC(app, dep)
		defer func(grpcServer *grpc.Server) {
			grpcServer.Stop()
		}(grpcServer)

		<-ctx.Done()
	case consumerMode:
		consumeKafkaMessages(ctx, app, dep)
	default:
		httpServer := serveHTTP(app, dep)
		defer func(httpServer *http.Server) {
			_ = httpServer.Close()
		}(httpServer)

		grpcServer := serveGRPC(app, dep)
		defer func(grpcServer *grpc.Server) {
			grpcServer.Stop()
		}(grpcServer)

		consumeKafkaMessages(ctx, app, dep)

		<-ctx.Done()
	}
}
func serveGRPC(app *application.App, dep *application.Dependency) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", app.Config.Const.GRPCPort))
	if err != nil {
		log.Fatalf("cannot listen on %v\n", app.Config.Const.GRPCPort)
	}

	grpcServer := grpc.NewServer()

	stockGRPCController := &controller.StockGRPCController{
		StockService: dep.StockService,
	}

	pb.RegisterStockServerServer(grpcServer, stockGRPCController)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc serve error: %s\n", err)
		}
	}()

	log.Printf("grpc running at :%v", app.Config.Const.GRPCPort)
	return grpcServer
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
			log.Fatalf("Failed to serve http server")
		}
	}(s)

	return s
}

func consumeKafkaMessages(ctx context.Context, app *application.App, dep *application.Dependency) {
	topics := app.Config.Kafka.ConsumerTopics
	consumerHandler := &controller.ConsumerHandler{
		TransactionFeed: dep.TransactionFeedService,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(t string, h func([]byte) (bool, error)) {
		log.Printf("creating consumer for topic %s", t)
		kafkaListener(ctx, app, t, h)
		wg.Done()
	}(topics["transaction"], consumerHandler.Transaction)
}

func kafkaListener(ctx context.Context, app *application.App, topic string, messageHandler func([]byte) (bool, error)) {
	log.Printf("Kafka Listener(%s) START", topic)

loop:
	for {
		select {
		case <-ctx.Done():
			log.Printf("consumer stop signal %s", topic)
			break loop
		default:
			msg, err := app.Consumer.Consume(ctx, topic)
			if err != nil {
				if err == io.EOF || err == context.Canceled {
					log.Printf("shutting down kafka consumers %s", topic)
					break loop
				}

				log.Printf("error when consuming kafka message %s", topic)
				continue
			} else if msg.Value == nil {
				// no message, no error. skip
				continue
			}
			// Process the message.
			log.Printf("kafka message consumed %s[%d]%d", topic, msg.Partition, msg.Offset)

			if _, err := messageHandler(msg.Value.([]byte)); err != nil {
				log.Printf("error processing message %s[%d]%d", topic, msg.Partition, msg.Offset)
			}
		}
	}

	_ = app.Consumer.Close()
	log.Printf("Listener (%s) STOP", topic)
}
