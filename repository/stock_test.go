package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"stockbit-challenge/adapter/redis/mock"
	"stockbit-challenge/model"
)

func TestStockRepository_SetStockInfo(t *testing.T) {
	var (
		mockCtrl  = gomock.NewController(t)
		mockRedis = mock.NewMockIRedis(mockCtrl)

		stock = model.Stock{
			Code: "BBCA",
		}

		key       = fmt.Sprintf(keyStock, stock.Code)
		stockRepo = StockRepository{
			Redis: mockRedis,
		}
	)

	t.Run("positive", func(t *testing.T) {
		mockRedis.EXPECT().SetAsBytes(key, stock).Return(nil).Times(1)
		err := stockRepo.SetStockInfo(stock)
		assert.Nil(t, err)
	})

	t.Run("negative", func(t *testing.T) {
		mockRedis.EXPECT().SetAsBytes(key, stock).Return(errors.New("wow error")).Times(1)
		err := stockRepo.SetStockInfo(stock)
		assert.EqualError(t, err, "Err set stock info: wow error")
	})
}

func TestStockRepository_GetStockInfo(t *testing.T) {
	var (
		mockCtrl  = gomock.NewController(t)
		mockRedis = mock.NewMockIRedis(mockCtrl)

		stock = model.Stock{
			Code: "BBCA",
		}

		emptyStock = &model.Stock{}

		key       = fmt.Sprintf(keyStock, stock.Code)
		stockRepo = StockRepository{
			Redis: mockRedis,
		}
	)

	t.Run("positive", func(t *testing.T) {
		mockRedis.EXPECT().GetAndParseBytes(key, emptyStock).Return(nil).Times(1)
		res, err := stockRepo.GetStockInfo(stock.Code)
		assert.Empty(t, res)
		assert.Nil(t, err)
	})

	t.Run("negative", func(t *testing.T) {
		mockRedis.EXPECT().GetAndParseBytes(key, emptyStock).Return(errors.New("error")).Times(1)
		res, err := stockRepo.GetStockInfo(stock.Code)
		assert.Empty(t, res)
		assert.EqualError(t, err, "error")
	})
}
