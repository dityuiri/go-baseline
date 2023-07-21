package repository

import (
	"errors"
	"fmt"

	"stockbit-challenge/adapter/redis"
	"stockbit-challenge/model"
)

//go:generate mockgen -package=repository_mock -destination=../mock/repository/stock.go . IStockRepository

type (
	IStockRepository interface {
		SetStockInfo(stock model.Stock) error
		GetStockInfo(stockCode string) (*model.Stock, error)
	}

	StockRepository struct {
		Redis redis.IRedis
	}
)

const (
	keyStock = "stock:code:%s"
)

func (r *StockRepository) SetStockInfo(stock model.Stock) error {
	var key = fmt.Sprintf(keyStock, stock.Code)

	err := r.Redis.SetAsBytes(key, stock)
	if err != nil {
		return errors.New(fmt.Sprintf("Err set stock info: %s", err))
	}

	return err
}

func (r *StockRepository) GetStockInfo(stockCode string) (*model.Stock, error) {
	var (
		key    = fmt.Sprintf(keyStock, stockCode)
		result = &model.Stock{}
	)

	err := r.Redis.GetAndParseBytes(key, result)
	return result, err
}