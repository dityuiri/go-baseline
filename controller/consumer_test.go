package controller

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dityuiri/go-baseline/adapter/logger/mock"
	"github.com/dityuiri/go-baseline/common"
	serviceMock "github.com/dityuiri/go-baseline/mock/service"
	"github.com/dityuiri/go-baseline/model"
)

func TestConsumer_Placeholder(t *testing.T) {
	var (
		mockCtrl                   = gomock.NewController(t)
		mockLogger                 = mock.NewMockILogger(mockCtrl)
		mockPlaceholderFeedService = serviceMock.NewMockIPlaceholderFeedService(mockCtrl)

		ctx    = context.Background()
		msg, _ = common.JsonMarshal(model.PlaceholderMessage{})

		consumer = ConsumerHandler{
			Logger:                 mockLogger,
			PlaceholderFeedService: mockPlaceholderFeedService,
		}
	)

	defer mockCtrl.Finish()

	t.Run("positive", func(t *testing.T) {
		mockPlaceholderFeedService.EXPECT().PlaceholderRecorded(ctx, gomock.Any()).Return(true, nil)

		res, err := consumer.Placeholder(msg)
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("unmarshal failed", func(t *testing.T) {
		// Patching the unmarshal method
		jsonUnmarshal := json.Unmarshal
		common.JsonUnmarshal = func(data []byte, v any) error {
			return errors.New("error")
		}

		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		res, err := consumer.Placeholder(msg)
		assert.EqualError(t, err, "error")
		assert.True(t, res)

		common.JsonUnmarshal = jsonUnmarshal
	})
}
