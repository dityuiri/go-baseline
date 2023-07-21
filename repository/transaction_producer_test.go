package repository

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"stockbit-challenge/adapter/kafka/producer/mock"
	"stockbit-challenge/config"
	"stockbit-challenge/model"
)

func TestTransactionProducer_ProduceTrx(t *testing.T) {
	var (
		mockCtrl     = gomock.NewController(t)
		mockProducer = mock.NewMockIProducer(mockCtrl)

		transaction = model.Transaction{
			StockCode: "BBCA",
			OrderBook: 81239126391,
		}

		expectedErr = errors.New("huuuu")

		producer = TransactionProducer{
			Producer: mockProducer,
			KafkaConfig: &config.Kafka{
				ProducerTopics: map[string]string{
					"transaction":     "transaction",
					"transaction_dlq": "transaction_dlq",
				},
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		mockProducer.EXPECT().Produce("transaction", gomock.Any()).Return(nil)
		err := producer.ProduceTrx(transaction)
		assert.Nil(t, err)
	})

	t.Run("error produce", func(t *testing.T) {
		mockProducer.EXPECT().Produce("transaction", gomock.Any()).Return(expectedErr)
		err := producer.ProduceTrx(transaction)
		assert.EqualError(t, expectedErr, err.Error())
	})
}

func TestTransactionProducer_ProduceTrxDLQ(t *testing.T) {
	var (
		mockCtrl     = gomock.NewController(t)
		mockProducer = mock.NewMockIProducer(mockCtrl)

		expectedErr = errors.New("huuuu")
		transaction = model.Transaction{
			StockCode: "BBCA",
			OrderBook: 81239126391,
			Error:     expectedErr.Error(),
		}

		producer = TransactionProducer{
			Producer: mockProducer,
			KafkaConfig: &config.Kafka{
				ProducerTopics: map[string]string{
					"transaction":     "transaction",
					"transaction_dlq": "transaction_dlq",
				},
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		mockProducer.EXPECT().Produce("transaction_dlq", gomock.Any()).Return(nil)
		err := producer.ProduceTrxDLQ(transaction, expectedErr)
		assert.Nil(t, err)
	})

	t.Run("error produce", func(t *testing.T) {
		mockProducer.EXPECT().Produce("transaction_dlq", gomock.Any()).Return(expectedErr)
		err := producer.ProduceTrxDLQ(transaction, expectedErr)
		assert.EqualError(t, expectedErr, err.Error())
	})
}
