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

	// Create a WaitGroup to track the goroutines.
	var wg sync.WaitGroup

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
		httpServer := serveHTTP(ctx, app, dep, &wg)
		defer func(httpServer *http.Server) {
			_ = httpServer.Close()
		}(httpServer)

		grpcServer := serveGRPC(ctx, app, dep, &wg)
		defer func(grpcServer *grpc.Server) {
			grpcServer.GracefulStop()
		}(grpcServer)

		wg.Wait()
	case consumerMode:
		consumeKafkaMessages(ctx, app, dep, &wg)
	default:
		httpServer := serveHTTP(ctx, app, dep, &wg)
		defer func(httpServer *http.Server) {
			_ = httpServer.Close()
		}(httpServer)

		grpcServer := serveGRPC(ctx, app, dep, &wg)
		defer func(grpcServer *grpc.Server) {
			grpcServer.GracefulStop()
		}(grpcServer)

		consumeKafkaMessages(ctx, app, dep, &wg)
		wg.Wait()
	}
}
func serveGRPC(ctx context.Context, app *application.App, dep *application.Dependency, wg *sync.WaitGroup) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", app.Config.Const.GRPCPort))
	if err != nil {
		log.Fatalf("cannot listen on %v\n", app.Config.Const.GRPCPort)
	}

	grpcServer := grpc.NewServer()

	stockGRPCController := &controller.StockGRPCController{
		StockService: dep.StockService,
	}

	pb.RegisterStockServerServer(grpcServer, stockGRPCController)

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Start the gRPC server
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("grpc serve error: %s\n", err)
			}
		}()
		<-ctx.Done() // Wait for the signal from the main context to shut down

		grpcServer.GracefulStop()
	}()

	log.Printf("grpc running at :%v", app.Config.Const.GRPCPort)
	return grpcServer
}

func serveHTTP(ctx context.Context, app *application.App, dep *application.Dependency, wg *sync.WaitGroup) *http.Server {
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

	// Increment the WaitGroup counter.
	wg.Add(1)
	go func(hs *http.Server) {
		defer wg.Done()

		// Start the HTTP server
		go func() {
			err := hs.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("failed to serve http server: %v", err)
			}
		}()

		<-ctx.Done() // Wait for the signal from the main context to shut down

		// Create a shutdown context with a timeout to allow existing connections to finish
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := hs.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("failed to gracefully shut down http server: %v", err)
		}
	}(s)

	return s
}

func consumeKafkaMessages(ctx context.Context, app *application.App, dep *application.Dependency, wg *sync.WaitGroup) {
	topics := app.Config.Kafka.ConsumerTopics
	consumerHandler := &controller.ConsumerHandler{
		TransactionFeed: dep.TransactionFeedService,
	}

	wg.Add(1)

	go func(t string, h func([]byte) (bool, error)) {
		defer wg.Done()
		log.Printf("creating consumer for topic %s", t)
		kafkaListener(ctx, app, t, h)
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
