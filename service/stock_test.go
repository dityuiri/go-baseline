package service

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	repositoryMock "stockbit-challenge/mock/repository"
	"stockbit-challenge/model"
)

var _ = Describe("StockService", func() {
	var (
		mockCtrl      *gomock.Controller
		mockStockRepo *repositoryMock.MockIStockRepository

		stockCode string
		service   StockService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockStockRepo = repositoryMock.NewMockIStockRepository(mockCtrl)

		stockCode = "BBCA"
		service = StockService{
			StockRepository: mockStockRepo,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetStockSummary", func() {
		Context("Positive", func() {
			It("should return an empty stock summary and nil error on success", func() {
				mockStockRepo.EXPECT().GetStockInfo(stockCode).Return(&model.Stock{}, nil)
				res, err := service.GetStockSummary(stockCode)
				Expect(res.Code).To(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Negative", func() {
			It("should return an empty stock summary and the expected error on failure", func() {
				expectedErr := errors.New("error")
				mockStockRepo.EXPECT().GetStockInfo(stockCode).Return(&model.Stock{}, expectedErr)
				res, err := service.GetStockSummary(stockCode)
				Expect(res.Code).To(BeEmpty())
				Expect(err).To(MatchError(expectedErr))
			})
		})
	})
})
