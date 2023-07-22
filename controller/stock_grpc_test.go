package controller

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	serviceMock "stockbit-challenge/mock/service"
	"stockbit-challenge/model"
	pb "stockbit-challenge/proto/proto-golang/stock"
)

func setupTestServer(t FullGinkgoTInterface, server *StockGRPCController) (*grpc.Server, *grpc.ClientConn) {
	// Create a new gRPC server
	s := grpc.NewServer()

	// Initialize your Server and register it with the gRPC server
	pb.RegisterStockServerServer(s, server)

	// Start a listener on a random port
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}

	// Start the gRPC server in a separate goroutine
	go func() {
		if err = s.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
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

var _ = Describe("StockGRPCController", func() {
	var (
		mockCtrl         *gomock.Controller
		mockStockService *serviceMock.MockIStockService

		ctx        context.Context
		controller *StockGRPCController

		server *grpc.Server
		conn   *grpc.ClientConn
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockStockService = serviceMock.NewMockIStockService(mockCtrl)

		ctx = context.Background()

		controller = &StockGRPCController{
			StockService: mockStockService,
		}

		server, conn = setupTestServer(GinkgoT(), controller)
	})

	AfterEach(func() {
		server.Stop()
		_ = conn.Close()
		mockCtrl.Finish()
	})

	Context("Without code", func() {
		It("should return an invalid argument error", func() {
			client := pb.NewStockServerClient(conn)
			stockRequest := &pb.Stock{}

			response, err := client.GetSummary(ctx, stockRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.GetError().GetCode()).To(Equal(int32(codes.InvalidArgument)))
		})
	})

	Context("Stock not found", func() {
		It("should return a not found error", func() {
			client := pb.NewStockServerClient(conn)
			stockRequest := &pb.Stock{Code: "BBCA"}

			mockStockService.EXPECT().GetStockSummary("BBCA").Return(&model.Stock{}, redis.Nil)
			response, err := client.GetSummary(ctx, stockRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.GetError().GetCode()).To(Equal(int32(codes.NotFound)))
		})
	})

	Context("Internal server error", func() {
		It("should return an internal server error", func() {
			client := pb.NewStockServerClient(conn)
			stockRequest := &pb.Stock{Code: "BBCA"}

			mockStockService.EXPECT().GetStockSummary("BBCA").Return(&model.Stock{}, errors.New("error"))
			response, err := client.GetSummary(ctx, stockRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.GetError().GetCode()).To(Equal(int32(codes.Internal)))
		})
	})

	Context("Success", func() {
		It("should return the stock summary", func() {
			client := pb.NewStockServerClient(conn)
			stockRequest := &pb.Stock{Code: "BBCA"}

			mockStockService.EXPECT().GetStockSummary("BBCA").Return(&model.Stock{Code: "BBCA"}, nil)
			response, err := client.GetSummary(ctx, stockRequest)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.GetSummary().GetCode()).To(Equal("BBCA"))
		})
	})
})
