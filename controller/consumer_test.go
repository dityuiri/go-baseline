package controller

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/golang/mock/gomock"

	serviceMock "stockbit-challenge/mock/service"
	"stockbit-challenge/model"
)

func TestConsumer_Transaction(t *testing.T) {
	var (
		mockCtrl    = gomock.NewController(t)
		mockTrxFeed = serviceMock.NewMockITransactionFeedService(mockCtrl)

		transaction = model.Transaction{
			Type: "A",
		}

		consumer = ConsumerHandler{
			TransactionFeed: mockTrxFeed,
		}
	)

	defer mockCtrl.Finish()

	t.Run("negative invalid message", func(t *testing.T) {
		var message = []byte("kechawww")

		res, err := consumer.Transaction(message)
		assert.NotNil(t, err)
		assert.True(t, res)
	})

	t.Run("feed returns error", func(t *testing.T) {
		request, _ := json.Marshal(&transaction)
		mockTrxFeed.EXPECT().TransactionRecorded(gomock.Any()).Return(false, errors.New("error"))
		res, err := consumer.Transaction(request)
		assert.False(t, res)
		assert.EqualError(t, err, "error")
	})

	t.Run("positive", func(t *testing.T) {
		request, _ := json.Marshal(&transaction)
		mockTrxFeed.EXPECT().TransactionRecorded(gomock.Any()).Return(true, nil)
		res, err := consumer.Transaction(request)
		assert.True(t, res)
		assert.Nil(t, err)
	})

}
