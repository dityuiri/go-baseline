package service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"

	"stockbit-challenge/model"
	"stockbit-challenge/repository"
)

//go:generate mockgen -package=service_mock -destination=../mock/service/transaction_feed.go . ITransactionFeedService

type (
	ITransactionFeedService interface {
		TransactionRecorded(trx *model.Transaction) (bool, error)
		ProduceTransaction(buff bytes.Buffer) error
	}

	TransactionFeedService struct {
		StockRepository     repository.IStockRepository
		TransactionProducer repository.ITransactionProducer
	}
)

func (s *TransactionFeedService) TransactionRecorded(trx *model.Transaction) (bool, error) {
	var err error

	// Validate transaction
	isValid := s.isValidTransaction(trx)
	if !isValid {
		return true, err
	}

	// Get stock info
	var stockInfo *model.Stock
	stockInfo, err = s.StockRepository.GetStockInfo(trx.StockCode)
	if err != nil {
		if err != redis.Nil {
			fmt.Println(fmt.Sprintf("Error getting stock info %s", err.Error()))
			s.produceDLQ(trx, err) // Produce DLQ and return false for systemic error
			return false, err
		}

		// Assumption: First occurrence of transaction must be previous price not accountable for OHLC
		// Added some validation here to make sure the assumption is strictly followed
		if !trx.IsPreviousPrice() || trx.IsAccountable() {
			fmt.Println(fmt.Sprintf("First occurrence must be previous price and not accountable type"))
			return true, nil // Skip event
		}

		// Create new stock info
		stockInfo = &model.Stock{
			Code:          trx.StockCode,
			PreviousPrice: trx.Price,
		}
	}

	// Case found in redis
	if err == nil {
		// Skip if it's previous price
		if trx.IsPreviousPrice() {
			fmt.Println(fmt.Sprintf("Previous price transaction while stock already exists"))
			return true, nil
		}

		// Skip type "A"
		if !trx.IsAccountable() {
			return true, nil
		}

		// Only update when transaction is accountable
		// First transaction of the day case
		if stockInfo.OpenPrice == 0 {
			stockInfo.OpenPrice = trx.Price
		}

		// New highest price found case
		if trx.Price > stockInfo.HighestPrice {
			stockInfo.HighestPrice = trx.Price
		}

		// New lowest price found case
		if trx.Price < stockInfo.LowestPrice || stockInfo.LowestPrice == 0 {
			stockInfo.LowestPrice = trx.Price
		}

		// Always update close price
		stockInfo.ClosePrice = trx.Price

		stockInfo.Volume += trx.Quantity
		stockInfo.Value += trx.Quantity * trx.Price
		stockInfo.AveragePrice = float64(stockInfo.Value / stockInfo.Volume)
	}

	// Set new stock info into the redis
	err = s.StockRepository.SetStockInfo(*stockInfo)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error setting value stock info %s", err.Error()))
		s.produceDLQ(trx, err) // Produce DLQ and return false for systemic error
		return false, err
	}

	return true, err
}

func (s *TransactionFeedService) isValidTransaction(trx *model.Transaction) bool {
	var isValidMandatoryValues = true

	if trx.StockCode == "" || trx.OrderBook == 0 || trx.Price == 0 {
		fmt.Println("Missing mandatory values")
		isValidMandatoryValues = false
	}

	return s.isValidTrxType(trx.Type) && isValidMandatoryValues
}

func (s *TransactionFeedService) isValidTrxType(trxType string) bool {
	for _, item := range model.TransactionTypes {
		if item == trxType {
			return true
		}
	}

	fmt.Println("Type not found")
	return false
}

func (s *TransactionFeedService) produceDLQ(transaction *model.Transaction, err error) {
	go func() {
		_ = s.TransactionProducer.ProduceTrxDLQ(*transaction, err)
	}()
}

func (s *TransactionFeedService) ProduceTransaction(buff bytes.Buffer) error {
	lines := bytes.Split(buff.Bytes(), []byte("\n"))
	for _, rawTx := range lines {
		var rawTransaction model.RawTransaction

		if err := json.Unmarshal(rawTx, &rawTransaction); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return err
		}

		transaction, err := rawTransaction.ToTransaction()
		if err != nil {
			fmt.Println("Error converting raw transaction:", err)
			return err
		}

		// Since it's not mandatory to have this helper for producing trx, skipping the error line
		go func() {
			err = s.TransactionProducer.ProduceTrx(transaction)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error producing transaction for %+v", transaction))
			}
		}()
	}

	return nil
}
