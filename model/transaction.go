package model

import (
	"errors"
	"strconv"
)

type (
	RawTransaction struct {
		Type             string `json:"type"`
		OrderNumber      string `json:"order_number"`
		OrderVerb        string `json:"order_verb"`
		Quantity         string `json:"quantity"`
		OrderBook        string `json:"order_book"`
		Price            string `json:"price"`
		StockCode        string `json:"stock_code"`
		ExecutedQuantity string `json:"executed_quantity"`
		ExecutedPrice    string `json:"executed_price"`
	}

	Transaction struct {
		Type        string `json:"type"`
		OrderNumber string `json:"order_number"`
		OrderVerb   string `json:"order_verb"`
		Quantity    int64  `json:"quantity"`
		OrderBook   int64  `json:"order_book"`
		Price       int64  `json:"price"`
		StockCode   string `json:"stock_code"`
		Error       string `json:"omitempty,error"`
	}
)

func (r *RawTransaction) ToTransaction() (Transaction, error) {
	var (
		transaction Transaction
		qty         = r.Quantity
		price       = r.Price
	)

	if r.Type != "A" {
		qty = r.ExecutedQuantity
		price = r.ExecutedPrice
	}

	finalQty, err := strconv.ParseInt(qty, 10, 64)
	if err != nil {
		return transaction, errors.New("invalid quantity format")
	}

	finalPrice, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		return transaction, errors.New("invalid price format")
	}

	orderBook, err := strconv.ParseInt(r.OrderBook, 10, 64)
	if err != nil {
		return transaction, errors.New("invalid order book format")
	}

	transaction = Transaction{
		Type:        r.Type,
		OrderBook:   orderBook,
		OrderNumber: r.OrderNumber,
		Quantity:    finalQty,
		Price:       finalPrice,
		StockCode:   r.StockCode,
	}

	return transaction, err
}

func (t *Transaction) IsPreviousPrice() bool {
	return t.Quantity == 0
}

func (t *Transaction) IsAccountable() bool {
	return t.Type == "E" || t.Type == "P"
}
