package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawTransaction_ToTransaction(t *testing.T) {
	var (
		rawTypeE = RawTransaction{
			Type:             "E",
			OrderNumber:      "12831289312389",
			OrderVerb:        "S",
			Quantity:         "13",
			Price:            "20000",
			StockCode:        "BBCA",
			ExecutionPrice:   "20000",
			ExecutedQuantity: "13",
			OrderBook:        "123",
		}
	)

	t.Run("positive", func(t *testing.T) {
		trx, err := rawTypeE.ToTransaction()
		assert.NotEmpty(t, trx)
		assert.Nil(t, err)
	})

	t.Run("invalid quantity", func(t *testing.T) {
		rawTypeInvalid := RawTransaction{
			Type:             "E",
			OrderNumber:      "12831289312389",
			OrderVerb:        "S",
			Quantity:         "13",
			Price:            "20000",
			StockCode:        "BBCA",
			ExecutionPrice:   "20000",
			ExecutedQuantity: "HAAAA",
		}

		trx, err := rawTypeInvalid.ToTransaction()
		assert.Empty(t, trx)
		assert.Error(t, err)
	})

	t.Run("invalid price", func(t *testing.T) {
		rawTypeInvalid := RawTransaction{
			Type:             "E",
			OrderNumber:      "12831289312389",
			OrderVerb:        "S",
			Quantity:         "13",
			Price:            "20000",
			StockCode:        "BBCA",
			ExecutionPrice:   "Sachi",
			ExecutedQuantity: "13",
		}

		trx, err := rawTypeInvalid.ToTransaction()
		assert.Empty(t, trx)
		assert.Error(t, err)
	})

	t.Run("invalid order book", func(t *testing.T) {
		rawTypeInvalid := RawTransaction{
			Type:             "E",
			OrderNumber:      "12831289312389",
			OrderVerb:        "S",
			Quantity:         "13",
			Price:            "20000",
			StockCode:        "BBCA",
			ExecutionPrice:   "20000",
			ExecutedQuantity: "13",
			OrderBook:        "JASDASD",
		}

		trx, err := rawTypeInvalid.ToTransaction()
		assert.Empty(t, trx)
		assert.Error(t, err)
	})
}

func TestTransaction_IsAccountable(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var trx = Transaction{Type: "E"}
		res := trx.IsAccountable()
		assert.True(t, res)
	})

	t.Run("false", func(t *testing.T) {
		var trx = Transaction{Type: "BLAU"}
		res := trx.IsAccountable()
		assert.False(t, res)
	})
}

func TestTransaction_IsPreviousPrice(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var trx = Transaction{Type: "E"}
		res := trx.IsPreviousPrice()
		assert.True(t, res)
	})

	t.Run("false", func(t *testing.T) {
		var trx = Transaction{Type: "BLAU"}
		res := trx.IsAccountable()
		assert.False(t, res)
	})
}
