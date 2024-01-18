package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	loggerMock "github.com/dityuiri/go-baseline/adapter/logger/mock"
	redisMock "github.com/dityuiri/go-baseline/adapter/redis/mock"
	"github.com/dityuiri/go-baseline/model"
)

func TestPlaceholderCache_SetPlaceholderInfo(t *testing.T) {
	var (
		mockCtrl   = gomock.NewController(t)
		mockRedis  = redisMock.NewMockIRedis(mockCtrl)
		mockLogger = loggerMock.NewMockILogger(mockCtrl)

		placeholderCache = PlaceholderCache{
			Redis:  mockRedis,
			Logger: mockLogger,
		}

		placeholderDTO = model.PlaceholderDTO{
			ID: uuid.New(),
		}

		ctx = context.Background()
		key = fmt.Sprintf(keyPlaceholder, placeholderDTO.ID.String())
	)

	t.Run("return ok", func(t *testing.T) {
		mockRedis.EXPECT().SetAsBytes(key, placeholderDTO).Return(nil).Times(1)

		err := placeholderCache.SetPlaceholderInfo(ctx, placeholderDTO)
		assert.Nil(t, err)
	})

	t.Run("return error", func(t *testing.T) {
		mockRedis.EXPECT().SetAsBytes(key, placeholderDTO).Return(errors.New("error")).Times(1)

		err := placeholderCache.SetPlaceholderInfo(ctx, placeholderDTO)
		assert.EqualError(t, err, "error")
	})
}

func TestPlaceholderCache_GetPlaceholderInfo(t *testing.T) {
	var (
		mockCtrl   = gomock.NewController(t)
		mockRedis  = redisMock.NewMockIRedis(mockCtrl)
		mockLogger = loggerMock.NewMockILogger(mockCtrl)

		placeholderCache = PlaceholderCache{
			Redis:  mockRedis,
			Logger: mockLogger,
		}

		placeholderDTO = model.PlaceholderDTO{
			ID: uuid.New(),
		}

		ctx = context.Background()
		key = fmt.Sprintf(keyPlaceholder, placeholderDTO.ID.String())
	)

	t.Run("return ok", func(t *testing.T) {
		mockRedis.EXPECT().GetAndParseBytes(key, gomock.Any()).Return(nil).Times(1)

		res, err := placeholderCache.GetPlaceholderInfo(ctx, placeholderDTO.ID.String())
		assert.Empty(t, res)
		assert.Nil(t, err)
	})

	t.Run("return error", func(t *testing.T) {
		mockRedis.EXPECT().GetAndParseBytes(key, gomock.Any()).Return(errors.New("error")).Times(1)

		res, err := placeholderCache.GetPlaceholderInfo(ctx, placeholderDTO.ID.String())
		assert.Empty(t, res)
		assert.EqualError(t, err, "error")
	})
}
