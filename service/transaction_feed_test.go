package service

import (
	"bytes"
	"context"
	"errors"
	"sync"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	repositoryMock "stockbit-challenge/mock/repository"
	"stockbit-challenge/model"
)

var _ = Describe("TransactionFeedService", func() {
	var (
		mockCtrl        *gomock.Controller
		mockStockRepo   *repositoryMock.MockIStockRepository
		mockTrxProducer *repositoryMock.MockITransactionProducer

		expectedErr   error
		existingStock *model.Stock
		wg            sync.WaitGroup
		service       TransactionFeedService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockStockRepo = repositoryMock.NewMockIStockRepository(mockCtrl)
		mockTrxProducer = repositoryMock.NewMockITransactionProducer(mockCtrl)

		expectedErr = errors.New("error")
		wg = sync.WaitGroup{}
		existingStock = &model.Stock{
			Code:          "BBAW",
			PreviousPrice: 120,
		}
		service = TransactionFeedService{
			StockRepository:     mockStockRepo,
			TransactionProducer: mockTrxProducer,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("TransactionRecorded", func() {
		Context("Skipping invalid transaction", func() {
			It("should return true and nil error", func() {
				trx := &model.Transaction{}
				success, err := service.TransactionRecorded(trx)
				Expect(success).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Get stock info return error", func() {
			It("should return false and the expected error", func() {
				initialTrx := &model.Transaction{
					Type:      "A",
					StockCode: "BBAW",
					OrderBook: 123123,
					Price:     700,
				}

				wg.Add(1)
				defer wg.Wait()

				mockStockRepo.EXPECT().GetStockInfo(initialTrx.StockCode).Return(&model.Stock{}, expectedErr).Times(1)
				mockTrxProducer.EXPECT().ProduceTrxDLQ(gomock.Any(), *initialTrx, expectedErr).Do(func(ctx context.Context, trx model.Transaction, err error) {
					wg.Done()
				}).Return(nil)

				success, err := service.TransactionRecorded(initialTrx)
				Expect(success).To(BeFalse())
				Expect(err).To(MatchError(expectedErr))
			})
		})

		Context("Invalid initial transaction", func() {
			It("should return true and nil error", func() {
				invalidInitialTrx := &model.Transaction{
					Type:      "E",
					StockCode: "BBAW",
					OrderBook: 123123,
					Quantity:  20,
					Price:     600,
				}

				mockStockRepo.EXPECT().GetStockInfo(invalidInitialTrx.StockCode).Return(&model.Stock{}, redis.Nil).Times(1)

				success, err := service.TransactionRecorded(invalidInitialTrx)
				Expect(success).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Set stock info return error", func() {
			It("should return false and the expected error", func() {
				initialTrx := &model.Transaction{
					Type:      "A",
					StockCode: "BBAW",
					OrderBook: 123123,
					Price:     700,
				}

				wg.Add(1)
				defer wg.Wait()

				mockStockRepo.EXPECT().GetStockInfo(initialTrx.StockCode).Return(&model.Stock{}, redis.Nil).Times(1)
				mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(expectedErr).Times(1)
				mockTrxProducer.EXPECT().ProduceTrxDLQ(gomock.Any(), *initialTrx, expectedErr).Do(func(ctx context.Context, trx model.Transaction, err error) {
					wg.Done()
				}).Return(nil)

				success, err := service.TransactionRecorded(initialTrx)
				Expect(success).To(BeFalse())
				Expect(err).To(MatchError(expectedErr))
			})
		})

		Context("Success - initial transaction", func() {
			It("should return true and nil error", func() {
				initialTrx := &model.Transaction{
					Type:      "A",
					StockCode: "BBAW",
					OrderBook: 123123,
					Price:     700,
				}

				mockStockRepo.EXPECT().GetStockInfo(initialTrx.StockCode).Return(&model.Stock{}, redis.Nil).Times(1)
				mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(nil).Times(1)

				success, err := service.TransactionRecorded(initialTrx)
				Expect(success).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Stock exists yet transaction is for previous price", func() {
			It("should return true and nil error", func() {
				invalidTrx := &model.Transaction{
					Type:      "A",
					Quantity:  0,
					Price:     52,
					OrderBook: 124124124,
					StockCode: "BBAW",
				}

				mockStockRepo.EXPECT().GetStockInfo(invalidTrx.StockCode).Return(existingStock, nil).Times(1)

				success, err := service.TransactionRecorded(invalidTrx)
				Expect(success).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Stock exists yet type is A", func() {
			It("should return true and nil error", func() {
				invalidTrx := &model.Transaction{
					Type:      "A",
					Quantity:  200,
					Price:     52,
					OrderBook: 124124124,
					StockCode: "BBAW",
				}

				mockStockRepo.EXPECT().GetStockInfo(invalidTrx.StockCode).Return(existingStock, nil).Times(1)

				success, err := service.TransactionRecorded(invalidTrx)
				Expect(success).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Success - stock exists - highest updated", func() {
			It("should return true and nil error", func() {
				trx := &model.Transaction{
					Type:      "E",
					Quantity:  200,
					Price:     52,
					OrderBook: 124124124,
					StockCode: "BBAW",
				}

				mockStockRepo.EXPECT().GetStockInfo(trx.StockCode).Return(existingStock, nil).Times(1)
				mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(nil).Times(1)

				success, err := service.TransactionRecorded(trx)
				Expect(success).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Success - stock exists - lowest updated", func() {
			It("should return true and nil error", func() {
				trx := &model.Transaction{
					Type:      "E",
					Quantity:  200,
					Price:     52,
					OrderBook: 124124124,
					StockCode: "BBAW",
				}

				newExistingStock := &model.Stock{
					Code:          "BBAW",
					PreviousPrice: 120,
					LowestPrice:   400,
				}

				mockStockRepo.EXPECT().GetStockInfo(trx.StockCode).Return(newExistingStock, nil).Times(1)
				mockStockRepo.EXPECT().SetStockInfo(gomock.Any()).Return(nil).Times(1)

				success, err := service.TransactionRecorded(trx)
				Expect(success).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("ProduceTransaction", func() {
		Context("Unmarshal error", func() {
			It("should be ok", func() {
				invalidNdjsonData := []byte(`\n`)
				buff := bytes.NewBuffer(invalidNdjsonData)

				wg.Add(1)

				go func() {
					defer wg.Done()
					err := service.ProduceTransaction(*buff)
					Expect(err).NotTo(HaveOccurred())
				}()

				wg.Wait()
			})
		})

		Context("Invalid raw transaction data", func() {
			It("should be ok", func() {
				invalidNdjsonData := []byte(`{"type":"A","order_book":"35","price":"haha","stock_code":"UNVR"}`)
				buff := bytes.NewBuffer(invalidNdjsonData)

				wg.Add(1)

				go func() {
					defer wg.Done()
					err := service.ProduceTransaction(*buff)
					Expect(err).NotTo(HaveOccurred())
				}()

				wg.Wait()
			})
		})

		Context("Positive - all ok", func() {
			It("should not return an error", func() {

				ndjsonData := []byte(`{"type":"A","order_book":"35","price":"4540","stock_code":"UNVR"}`)
				buff := bytes.NewBuffer(ndjsonData)

				wg.Add(1)

				mockTrxProducer.EXPECT().ProduceTrx(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, transaction model.Transaction) {
					defer wg.Done()
				}).Return(nil).Times(1)

				go func() {
					err := service.ProduceTransaction(*buff)
					Expect(err).ToNot(HaveOccurred())
				}()

				wg.Wait()
			})
		})

		Context("Positive - produce error", func() {
			It("should not return an error", func() {
				ndjsonData := []byte(`{"type":"A","order_book":"35","price":"4540","stock_code":"UNVR"}`)
				buff := bytes.NewBuffer(ndjsonData)

				wg.Add(1)
				mockTrxProducer.EXPECT().ProduceTrx(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, transaction model.Transaction) {
					defer wg.Done()
				}).Return(errors.New("error")).Times(1)

				go func() {
					err := service.ProduceTransaction(*buff)
					Expect(err).ToNot(HaveOccurred())
				}()

				wg.Wait()
			})
		})
	})
})
