package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	repositoryMock "stockbit-challenge/mock/repository"
	"stockbit-challenge/model"
)

func TestStockService_GetStockSummary(t *testing.T) {
	var (
		mockCtrl      = gomock.NewController(t)
		mockStockRepo = repositoryMock.NewMockIStockRepository(mockCtrl)

		stockCode = "BBCA"
		service   = StockService{
			StockRepository: mockStockRepo,
		}
	)

	t.Run("positive", func(t *testing.T) {
		mockStockRepo.EXPECT().GetStockInfo(stockCode).Return(&model.Stock{}, nil)
		res, err := service.GetStockSummary(stockCode)
		assert.Empty(t, res)
		assert.Nil(t, err)
	})

	t.Run("negative", func(t *testing.T) {
		mockStockRepo.EXPECT().GetStockInfo(stockCode).Return(&model.Stock{}, errors.New("error"))
		res, err := service.GetStockSummary(stockCode)
		assert.Empty(t, res)
		assert.EqualError(t, err, "error")
	})
}
