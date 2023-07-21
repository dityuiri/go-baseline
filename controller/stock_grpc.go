package controller

import (
	"context"
	"math"

	"github.com/go-redis/redis"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	pb "stockbit-challenge/proto/proto-golang/stock"
	"stockbit-challenge/service"
)

type (
	IStockGRPCController interface {
		GetSummary(ctx context.Context, stockRequest *pb.Stock) (*pb.Response, error)
	}

	StockGRPCController struct {
		StockService service.IStockService
	}
)

func (g *StockGRPCController) GetSummary(ctx context.Context, stockRequest *pb.Stock) (*pb.Response, error) {
	if stockRequest.GetCode() == "" {
		return &pb.Response{
			Response: &pb.Response_Error{
				Error: &status.Status{
					Code:    int32(codes.InvalidArgument),
					Message: "code cannot be empty",
				},
			},
		}, nil
	}

	stockSummary, err := g.StockService.GetStockSummary(stockRequest.GetCode())
	if err != nil {
		var (
			code   = int32(codes.Internal)
			errMsg = err.Error()
		)

		if err == redis.Nil {
			code = int32(codes.NotFound)
			errMsg = "stock not found"
		}

		return &pb.Response{
			Response: &pb.Response_Error{
				Error: &status.Status{
					Code:    code,
					Message: errMsg,
				},
			},
		}, nil
	}

	// Return a successful response with the greeting message and int64 value
	return &pb.Response{
		Response: &pb.Response_Summary{
			Summary: &pb.StockSummary{
				Code:          stockSummary.Code,
				PreviousPrice: stockSummary.PreviousPrice,
				OpenPrice:     stockSummary.OpenPrice,
				HighestPrice:  stockSummary.HighestPrice,
				LowestPrice:   stockSummary.LowestPrice,
				ClosePrice:    stockSummary.ClosePrice,
				Volume:        stockSummary.Volume,
				Value:         stockSummary.Value,
				AveragePrice:  int64(math.Round(stockSummary.AveragePrice)),
			},
		},
	}, nil

}
