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
			ExecutedPrice:    "20000",
			ExecutedQuantity: "13",
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
			ExecutedPrice:    "20000",
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
			ExecutedPrice:    "Sachi",
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
			ExecutedPrice:    "20000",
			ExecutedQuantity: "13",
			OrderBook:        "JASDASD",
		}

		trx, err := rawTypeInvalid.ToTransaction()
		assert.Empty(t, trx)
		assert.Error(t, err)
	})
}
