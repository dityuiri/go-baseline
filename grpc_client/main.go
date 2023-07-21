package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"google.golang.org/grpc"

	pb "stockbit-challenge/proto/proto-golang/stock"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	client := pb.NewStockServerClient(conn)

	// Try calling with BBCA
	stockRequest := &pb.Stock{
		Code: "BBCA",
	}

	response, err := client.GetSummary(context.Background(), stockRequest)
	if err != nil {
		log.Fatalf("could not get summary: %v", err)
	}

	// Use the response from the server
	log.Printf("Stock Summary: %+v", response)
}
