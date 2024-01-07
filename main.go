package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/dityuiri/go-baseline/adapter/server"
	"github.com/dityuiri/go-baseline/application"
	"github.com/dityuiri/go-baseline/controller"
)

const (
	clientMode = "client"
	// consumerMode = "consumer"

	defaultTimeout = 5 * time.Second
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
	app, err := application.SetupApplication(ctx)
	if err != nil {
		panic(err)
	}

	defer app.Close()

	// Setup dependency injection
	dep := application.SetupDependency(app)

	// Goroutine to cancel when Os interrupt happens
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Received system interrupt. Stopping servers...")
		cancel()
	}()

	switch mode {
	case clientMode:
		var (
			httpServer = serveHTTP(app, dep)
		)
		if err := httpServer.Serve(); err != nil {
			panic(err)
		}

		<-app.Context.Done()
		_ = httpServer.Close()
	//case consumerMode:
	//	consumeKafkaMessages(app, dep)
	default: // All services run in this mode as default
		var (
			httpServer = serveHTTP(app, dep)
		)

		if err := httpServer.Serve(); err != nil {
			panic(err)
		}

		//consumeKafkaMessages(app, dep)

		<-app.Context.Done()
		_ = httpServer.Close()
	}
}

func serveHTTP(app *application.App, dep *application.Dependency) server.IServer {
	config := &server.Configuration{
		AppName: app.Config.AppName,
		Port:    app.Config.Const.HTTPPort,
	}

	httpServer := server.NewServer(app.Context, config)

	trxController := &controller.TransactionController{
		TransactionFeedService: dep.TransactionFeedService,
	}

	healthCheckController := &controller.HealthCheckController{
		HealthCheckService: dep.HealthCheckService,
	}
	httpServer.Get("/ping", healthCheckController.Ping)
	httpServer.Post("/publish/transaction", trxController.UploadTransactions)

	return httpServer
}

func withTimeout(timeout time.Duration) func(next http.Handler) http.Handler { //nolint
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			next.ServeHTTP(w, r)

			select {
			case <-ctx.Done():
				// TODO: if possible, respond with "Timeout" message
			default:
			}
		}

		return http.HandlerFunc(fn)
	}
}

// TODO: Uncomment if you want to use Kafka
//func consumeKafkaMessages(ctx context.Context, app *application.App, dep *application.Dependency, wg *sync.WaitGroup) {
//	topics := app.Config.Kafka.ConsumerTopics
//	consumerHandler := &controller.ConsumerHandler{
//		TransactionFeed: dep.TransactionFeedService,
//	}
//
//	wg.Add(1)
//
//	go func(t string, h func([]byte) (bool, error)) {
//		defer wg.Done()
//		log.Printf("creating consumer for topic %s", t)
//		kafkaListener(ctx, app, t, h)
//	}(topics["transaction"], consumerHandler.Transaction)
//}

//func kafkaListener(ctx context.Context, app *application.App, topic string, messageHandler func([]byte) (bool, error)) {
//	log.Printf("Kafka Listener(%s) START", topic)
//
//loop:
//	for {
//		select {
//		case <-ctx.Done():
//			log.Printf("consumer stop signal %s", topic)
//			break loop
//		default:
//			msg, err := app.Consumer.Consume(ctx, topic)
//			if err != nil {
//				if err == io.EOF || err == context.Canceled {
//					log.Printf("shutting down kafka consumers %s", topic)
//					break loop
//				}
//
//				log.Printf("error when consuming kafka message %s", topic)
//				continue
//			} else if msg.Value == nil {
//				// no message, no error. skip
//				continue
//			}
//			// Process the message.
//			log.Printf("kafka message consumed %s[%d]%d", topic, msg.Partition, msg.Offset)
//
//			if _, err := messageHandler(msg.Value.([]byte)); err != nil {
//				log.Printf("error processing message %s[%d]%d", topic, msg.Partition, msg.Offset)
//			}
//		}
//	}
//
//	_ = app.Consumer.Close()
//	log.Printf("Listener (%s) STOP", topic)
//}
