package service

import (
	"errors"
	"sync"
	"testing"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	repositoryMock "stockbit-challenge/mock/repository"
	"stockbit-challenge/model"
)

func TestTransactionFeedService_TransactionRecorded(t *testing.T) {
	var (
		mockCtrl        = gomock.NewController(t)
		mockStockRepo   = repositoryMock.NewMockIStockRepository(mockCtrl)
		mockTrxProducer = repositoryMock.NewMockITransactionProducer(mockCtrl)

		expectedErr = errors.New("error")
		wg          sync.WaitGroup
		initialTrx  = &model.Transaction{
			Type:      "A",
			StockCode: "BBAW",
			OrderBook: 123123,
			Price:     700,
		}

		existingStock = model.Stock{
			Code:          "BBAW",
			PreviousPrice: 120,
		}

		service = TransactionFeedService{
			StockRepository:     mockStockRepo,
			TransactionProducer: mockTrxProducer,
		}
	)

	defer mockCtrl.Finish()

	t.Run("skipping invalid transaction", func(t *testing.T) {
		trx := &model.Transaction{}
		success, err := service.TransactionRecorded(trx)
		assert.True(t, success)
		assert.Nil(t, err)
	})

	t.Run("get stock info return error", func(t *testing.T) {
		wg.Add(1)
		defer wg.Wait()

		mockStockRepo.EXPECT().GetStockInfo(initialTrx.StockCode).Return(model.Stock{}, expectedErr).Times(1)
		mockTrxProducer.EXPECT().ProduceTrxDLQ(*initialTrx, expectedErr).Do(func(trx model.Transaction, err error) {
			wg.Done()
		}).Return(nil)

		success, err := service.TransactionRecorded(initialTrx)
		assert.False(t, success)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("invalid initial transaction", func(t *testing.T) {
		invalidInitialTrx := &model.Transaction{
			Type:      "E",
			StockCode: "BBAW",
			OrderBook: 123123,
			Quantity:  20,
			Price:     600,
		}

		mockStockRepo.EXPECT().GetStockInfo(invalidInitialTrx.StockCode).Return(model.Stock{}, redis.Nil).Times(1)

		success, err := service.TransactionRecorded(invalidInitialTrx)
		assert.True(t, success)
		assert.Nil(t, err)
	})

	t.Run("set stock info return error", func(t *testing.T) {
		wg.Add(1)
		defer wg.Wait()

		mockStockRepo.EXPECT().GetStockInfo(initialTrx.StockCode).Return(model.Stock{}, redis.Nil).Times(1)
		mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(expectedErr).Times(1)
		mockTrxProducer.EXPECT().ProduceTrxDLQ(*initialTrx, expectedErr).Do(func(trx model.Transaction, err error) {
			wg.Done()
		}).Return(nil)

		success, err := service.TransactionRecorded(initialTrx)
		assert.False(t, success)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("success - initial transaction", func(t *testing.T) {
		mockStockRepo.EXPECT().GetStockInfo(initialTrx.StockCode).Return(model.Stock{}, redis.Nil).Times(1)
		mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(nil).Times(1)

		success, err := service.TransactionRecorded(initialTrx)
		assert.True(t, success)
		assert.Nil(t, err)
	})

	t.Run("stock exists yet transaction is for previous price", func(t *testing.T) {
		invalidTrx := &model.Transaction{
			Type:      "A",
			Quantity:  0,
			Price:     52,
			OrderBook: 124124124,
			StockCode: "BBAW",
		}

		mockStockRepo.EXPECT().GetStockInfo(invalidTrx.StockCode).Return(existingStock, nil).Times(1)

		success, err := service.TransactionRecorded(invalidTrx)
		assert.True(t, success)
		assert.Nil(t, err)
	})

	t.Run("stock exists yet type is A", func(t *testing.T) {
		invalidTrx := &model.Transaction{
			Type:      "A",
			Quantity:  200,
			Price:     52,
			OrderBook: 124124124,
			StockCode: "BBAW",
		}

		mockStockRepo.EXPECT().GetStockInfo(invalidTrx.StockCode).Return(existingStock, nil).Times(1)

		success, err := service.TransactionRecorded(invalidTrx)
		assert.True(t, success)
		assert.Nil(t, err)
	})

	t.Run("success - stock exists - highest updated", func(t *testing.T) {
		trx := &model.Transaction{
			Type:      "E",
			Quantity:  200,
			Price:     52,
			OrderBook: 124124124,
			StockCode: "BBAW",
		}

		mockStockRepo.EXPECT().GetStockInfo(trx.StockCode).Return(existingStock, nil).Times(1)
		mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(nil).Times(1)

		success, err := service.TransactionRecorded(trx)
		assert.True(t, success)
		assert.Nil(t, err)
	})

	t.Run("success - stock exists - lowest updated", func(t *testing.T) {
		trx := &model.Transaction{
			Type:      "E",
			Quantity:  200,
			Price:     52,
			OrderBook: 124124124,
			StockCode: "BBAW",
		}

		newExistingStock := model.Stock{
			Code:          "BBAW",
			PreviousPrice: 120,
			LowestPrice:   400,
		}

		mockStockRepo.EXPECT().GetStockInfo(trx.StockCode).Return(newExistingStock, nil).Times(1)
		mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(nil).Times(1)

		success, err := service.TransactionRecorded(trx)
		assert.True(t, success)
		assert.Nil(t, err)
	})

}
