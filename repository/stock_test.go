package repository

import (
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stockbit-challenge/adapter/redis/mock"
	"stockbit-challenge/model"
)

var _ = Describe("StockRepository", func() {
	var (
		mockCtrl  *gomock.Controller
		mockRedis *mock.MockIRedis
		stockRepo StockRepository

		stock = model.Stock{
			Code: "BBCA",
		}

		emptyStock = &model.Stock{}

		key = fmt.Sprintf(keyStock, stock.Code)
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRedis = mock.NewMockIRedis(mockCtrl)
		stockRepo = StockRepository{
			Redis: mockRedis,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("SetStockInfo", func() {
		It("should return nil error on success", func() {
			mockRedis.EXPECT().SetAsBytes(key, stock).Return(nil).Times(1)
			err := stockRepo.SetStockInfo(stock)
			Expect(err).To(BeNil())
		})

		It("should return error on set as bytes failure", func() {
			mockRedis.EXPECT().SetAsBytes(key, stock).Return(errors.New("error")).Times(1)
			err := stockRepo.SetStockInfo(stock)
			Expect(err).To(MatchError(errors.New("err set stock info: error")))
		})

	})

	Context("GetStockInfo", func() {
		It("should return empty result and nil error on success", func() {
			mockRedis.EXPECT().GetAndParseBytes(key, emptyStock).Return(nil).Times(1)
			res, err := stockRepo.GetStockInfo(stock.Code)
			Expect(res.Code).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return empty result and non-nil error on failure", func() {
			mockRedis.EXPECT().GetAndParseBytes(key, emptyStock).Return(errors.New("error")).Times(1)
			res, err := stockRepo.GetStockInfo(stock.Code)
			Expect(res.Code).To(BeEmpty())
			Expect(err).To(MatchError("error"))
		})
	})
})
