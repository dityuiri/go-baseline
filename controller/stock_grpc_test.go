package controller

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	serviceMock "stockbit-challenge/mock/service"
	"stockbit-challenge/model"
	pb "stockbit-challenge/proto/proto-golang/stock"
)

func setupTestServer(t *testing.T, server *StockGRPCController) (*grpc.Server, *grpc.ClientConn) {
	// Create a new gRPC server
	s := grpc.NewServer()

	// Initialize your Server and register it with the gRPC server
	pb.RegisterStockServerServer(s, server)

	// Start a listener on a random port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	// Start the gRPC server in a separate goroutine
	go func() {
		if err = s.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()

	// Create a gRPC client to connect to the test server
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	return s, conn
}

func TestStockGRPCController_GetSummary(t *testing.T) {
	var (
		mockCtrl         = gomock.NewController(t)
		mockStockService = serviceMock.NewMockIStockService(mockCtrl)

		ctx = context.Background()

		controller = &StockGRPCController{
			StockService: mockStockService,
		}

		server, conn = setupTestServer(t, controller)
	)

	defer server.Stop()
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	t.Run("without code", func(t *testing.T) {
		client := pb.NewStockServerClient(conn)
		stockRequest := &pb.Stock{}

		response, err := client.GetSummary(ctx, stockRequest)
		assert.Nil(t, err)
		assert.Equal(t, int32(codes.InvalidArgument), response.GetError().GetCode())
	})

	t.Run("stock not found", func(t *testing.T) {
		client := pb.NewStockServerClient(conn)
		stockRequest := &pb.Stock{Code: "BBCA"}

		mockStockService.EXPECT().GetStockSummary("BBCA").Return(&model.Stock{}, redis.Nil)
		response, err := client.GetSummary(ctx, stockRequest)
		assert.Nil(t, err)
		assert.Equal(t, int32(codes.NotFound), response.GetError().GetCode())
	})

	t.Run("internal server error", func(t *testing.T) {
		client := pb.NewStockServerClient(conn)
		stockRequest := &pb.Stock{Code: "BBCA"}

		mockStockService.EXPECT().GetStockSummary("BBCA").Return(&model.Stock{}, errors.New("error"))
		response, err := client.GetSummary(ctx, stockRequest)
		assert.Nil(t, err)
		assert.Equal(t, int32(codes.Internal), response.GetError().GetCode())
	})

	t.Run("success", func(t *testing.T) {
		client := pb.NewStockServerClient(conn)
		stockRequest := &pb.Stock{Code: "BBCA"}

		mockStockService.EXPECT().GetStockSummary("BBCA").Return(&model.Stock{Code: "BBCA"}, nil)
		response, err := client.GetSummary(ctx, stockRequest)
		assert.Nil(t, err)
		assert.Equal(t, "BBCA", response.GetSummary().GetCode())
	})
}
