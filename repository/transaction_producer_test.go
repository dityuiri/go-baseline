package repository

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mock2 "stockbit-challenge/adapter/kafka/producer/mock"

	"stockbit-challenge/config"
	"stockbit-challenge/model"
)

var _ = Describe("TransactionProducer", func() {
	var (
		mockCtrl     *gomock.Controller
		mockProducer *mock2.MockIProducer

		ctx         context.Context
		transaction model.Transaction
		expectedErr error

		producer TransactionProducer
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockProducer = mock2.NewMockIProducer(mockCtrl)

		ctx = context.Background()
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
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("ProduceTrx", func() {
		Context("Success", func() {
			It("should produce the transaction successfully", func() {
				mockProducer.EXPECT().Produce(ctx, "transaction", gomock.Any()).Return(nil)
				err := producer.ProduceTrx(ctx, transaction)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Error during produce", func() {
			It("should return the expected error", func() {
				mockProducer.EXPECT().Produce(ctx, "transaction", gomock.Any()).Return(expectedErr)
				err := producer.ProduceTrx(ctx, transaction)
				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Describe("ProduceTrxDLQ", func() {
		Context("Success", func() {
			It("should produce the transaction DLQ successfully", func() {
				mockProducer.EXPECT().Produce(ctx, "transaction_dlq", gomock.Any()).Return(nil)
				err := producer.ProduceTrxDLQ(ctx, transaction, expectedErr)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Error during produce", func() {
			It("should return the expected error", func() {
				mockProducer.EXPECT().Produce(ctx, "transaction_dlq", gomock.Any()).Return(expectedErr)
				err := producer.ProduceTrxDLQ(ctx, transaction, expectedErr)
				Expect(err).To(Equal(expectedErr))
			})
		})
	})
})
