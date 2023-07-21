package service

import (
	"stockbit-challenge/model"
	"stockbit-challenge/repository"
)

//go:generate mockgen -package=service_mock -destination=../mock/service/stock.go . IStockService

type (
	IStockService interface {
		GetStockSummary(stockCode string) (*model.Stock, error)
	}

	StockService struct {
		StockRepository repository.IStockRepository
	}
)

func (s *StockService) GetStockSummary(stockCode string) (*model.Stock, error) {
	return s.StockRepository.GetStockInfo(stockCode)
}
