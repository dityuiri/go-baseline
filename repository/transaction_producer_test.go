package repository

import (
	"context"
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

		ctx         = context.Background()
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
		mockProducer.EXPECT().Produce(ctx, "transaction", gomock.Any()).Return(nil)
		err := producer.ProduceTrx(ctx, transaction)
		assert.Nil(t, err)
	})

	t.Run("error produce", func(t *testing.T) {
		mockProducer.EXPECT().Produce(ctx, "transaction", gomock.Any()).Return(expectedErr)
		err := producer.ProduceTrx(ctx, transaction)
		assert.EqualError(t, expectedErr, err.Error())
	})
}

func TestTransactionProducer_ProduceTrxDLQ(t *testing.T) {
	var (
		mockCtrl     = gomock.NewController(t)
		mockProducer = mock.NewMockIProducer(mockCtrl)

		ctx         = context.Background()
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
		mockProducer.EXPECT().Produce(ctx, "transaction_dlq", gomock.Any()).Return(nil)
		err := producer.ProduceTrxDLQ(ctx, transaction, expectedErr)
		assert.Nil(t, err)
	})

	t.Run("error produce", func(t *testing.T) {
		mockProducer.EXPECT().Produce(ctx, "transaction_dlq", gomock.Any()).Return(expectedErr)
		err := producer.ProduceTrxDLQ(ctx, transaction, expectedErr)
		assert.EqualError(t, expectedErr, err.Error())
	})
}
