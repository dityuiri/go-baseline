package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	loggerMock "github.com/dityuiri/go-adapter/logger/mock"
	repositoryMock "github.com/dityuiri/go-baseline/mock/repository"
	"github.com/dityuiri/go-baseline/model"
)

func TestPlaceholderFeedService_PlaceholderRecorded(t *testing.T) {
	var (
		mockCtrl                = gomock.NewController(t)
		mockLogger              = loggerMock.NewMockILogger(mockCtrl)
		mockPlaceholderProducer = repositoryMock.NewMockIPlaceholderProducer(mockCtrl)

		placeholderFeedService = PlaceholderFeedService{
			Logger:              mockLogger,
			PlaceholderProducer: mockPlaceholderProducer,
		}

		ctx            = context.Background()
		placeholderMsg = model.PlaceholderMessage{}
	)

	t.Run("positive", func(t *testing.T) {
		mockPlaceholderProducer.EXPECT().ProducePlaceholderRecord(ctx, gomock.Any()).Return(nil).Times(1)

		isSuccess, err := placeholderFeedService.PlaceholderRecorded(ctx, placeholderMsg)
		assert.Nil(t, err)
		assert.True(t, isSuccess)
	})

	t.Run("producer returning error", func(t *testing.T) {
		mockPlaceholderProducer.EXPECT().ProducePlaceholderRecord(ctx, gomock.Any()).Return(errors.New("error")).Times(1)
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		isSuccess, err := placeholderFeedService.PlaceholderRecorded(ctx, placeholderMsg)
		assert.EqualError(t, err, "error")
		assert.False(t, isSuccess)
	})
}
