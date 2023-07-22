package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "stockbit-challenge/proto/proto-golang/stock"
)

func main() {
	var (
		args = os.Args[1:]
		code string
	)

	if len(args) > 0 {
		code = os.Args[1]
	}

	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	client := pb.NewStockServerClient(conn)

	// Code is from input
	stockRequest := &pb.Stock{
		Code: code,
	}

	response, err := client.GetSummary(context.Background(), stockRequest)
	if err != nil {
		log.Printf("could not get summary: %v", err)
	}

	// Use the response from the server
	log.Printf("Stock Summary: %+v", response)
}
