package controller

import (
	"encoding/json"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	serviceMock "github.com/dityuiri/go-baseline/mock/service"
	"github.com/dityuiri/go-baseline/model"
)

var _ = Describe("Consumer Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		mockTrxFeed *serviceMock.MockITransactionFeedService

		transaction model.Transaction

		consumer ConsumerHandler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTrxFeed = serviceMock.NewMockITransactionFeedService(mockCtrl)

		transaction = model.Transaction{
			Type: "A",
		}

		consumer = ConsumerHandler{
			TransactionFeed: mockTrxFeed,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Negative invalid message", func() {
		It("should return an error", func() {
			message := []byte("kechawww")

			res, err := consumer.Transaction(message)
			Expect(err).ToNot(BeNil())
			Expect(res).To(BeTrue())
		})
	})

	Context("Feed returns error", func() {
		It("should return an error", func() {
			request, _ := json.Marshal(&transaction)
			mockTrxFeed.EXPECT().TransactionRecorded(gomock.Any()).Return(false, errors.New("error"))

			res, err := consumer.Transaction(request)
			Expect(res).To(BeFalse())
			Expect(err).To(MatchError("error"))
		})
	})

	Context("Positive", func() {
		It("should not return an error", func() {
			request, _ := json.Marshal(&transaction)
			mockTrxFeed.EXPECT().TransactionRecorded(gomock.Any()).Return(true, nil)

			res, err := consumer.Transaction(request)
			Expect(res).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
