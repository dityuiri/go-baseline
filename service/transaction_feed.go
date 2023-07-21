package service

import (
	"stockbit-challenge/repository"

	"stockbit-challenge/model"
)

//go:generate mockgen -package=service_mock -destination=../mock/service/transaction_feed.go . ITransactionFeedService

type (
	ITransactionFeedService interface {
		TransactionRecorded(trx *model.Transaction) (bool, error)
	}

	TransactionFeedService struct {
		StockRepository     repository.IStockRepository
		TransactionProducer repository.ITransactionProducer
	}
)

func (s *TransactionFeedService) TransactionRecorded(trx *model.Transaction) (bool, error) {
	return false, nil
}
