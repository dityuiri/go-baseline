package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	loggerMock "github.com/dityuiri/go-adapter/logger/mock"
	"github.com/dityuiri/go-baseline/common"
	proxyMock "github.com/dityuiri/go-baseline/mock/proxy"
	repositoryMock "github.com/dityuiri/go-baseline/mock/repository"
	"github.com/dityuiri/go-baseline/model"
	"github.com/dityuiri/go-baseline/model/alpha"
)

func TestPlaceholderService_CreateNewPlaceholder(t *testing.T) {
	var (
		mockCtrl             = gomock.NewController(t)
		mockLogger           = loggerMock.NewMockILogger(mockCtrl)
		mockPlaceholderRepo  = repositoryMock.NewMockIPlaceholderRepository(mockCtrl)
		mockPlaceholderCache = repositoryMock.NewMockIPlaceholderCache(mockCtrl)
		mockAlphaProxy       = proxyMock.NewMockIAlphaProxy(mockCtrl)

		placeholderService = PlaceholderService{
			Logger:                mockLogger,
			PlaceholderRepository: mockPlaceholderRepo,
			PlaceholderCache:      mockPlaceholderCache,
			AlphaProxy:            mockAlphaProxy,
		}

		ctx                      = context.Background()
		placeholderCreateRequest = model.PlaceholderCreateRequest{}
	)

	t.Run("positive", func(t *testing.T) {
		mockPlaceholderRepo.EXPECT().InsertPlaceholder(ctx, nil, gomock.Any()).Return(nil)

		res, err := placeholderService.CreateNewPlaceholder(ctx, placeholderCreateRequest)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("insert placeholder returning error", func(t *testing.T) {
		mockPlaceholderRepo.EXPECT().InsertPlaceholder(ctx, nil, gomock.Any()).Return(errors.New("error"))
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		res, err := placeholderService.CreateNewPlaceholder(ctx, placeholderCreateRequest)
		assert.EqualError(t, err, "error")
		assert.Empty(t, res)
	})
}

func TestPlaceholderService_GetPlaceholder(t *testing.T) {
	var (
		mockCtrl             = gomock.NewController(t)
		mockLogger           = loggerMock.NewMockILogger(mockCtrl)
		mockPlaceholderRepo  = repositoryMock.NewMockIPlaceholderRepository(mockCtrl)
		mockPlaceholderCache = repositoryMock.NewMockIPlaceholderCache(mockCtrl)
		mockAlphaProxy       = proxyMock.NewMockIAlphaProxy(mockCtrl)

		placeholderService = PlaceholderService{
			Logger:                mockLogger,
			PlaceholderRepository: mockPlaceholderRepo,
			PlaceholderCache:      mockPlaceholderCache,
			AlphaProxy:            mockAlphaProxy,
		}

		ctx           = context.Background()
		placeholderID = uuid.New()
	)

	t.Run("positive - found the cache", func(t *testing.T) {
		mockPlaceholderCache.EXPECT().GetPlaceholderInfo(ctx, placeholderID.String()).Return(&model.PlaceholderDTO{}, nil).Times(1)
		mockAlphaProxy.EXPECT().GetPlaceholderStatus(ctx, gomock.Any()).Return(alpha.AlphaResponse{}, nil).Times(1)

		res, err := placeholderService.GetPlaceholder(ctx, placeholderID.String())
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("negative - found the cache - get status from proxy error", func(t *testing.T) {
		mockPlaceholderCache.EXPECT().GetPlaceholderInfo(ctx, placeholderID.String()).Return(&model.PlaceholderDTO{}, nil).Times(1)
		mockAlphaProxy.EXPECT().GetPlaceholderStatus(ctx, gomock.Any()).Return(alpha.AlphaResponse{}, errors.New("error")).Times(1)
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		res, err := placeholderService.GetPlaceholder(ctx, placeholderID.String())
		assert.EqualError(t, err, "error")
		assert.NotEmpty(t, res)
	})

	t.Run("negative - get from cache returning error", func(t *testing.T) {
		mockPlaceholderCache.EXPECT().GetPlaceholderInfo(ctx, placeholderID.String()).Return(&model.PlaceholderDTO{}, errors.New("error")).Times(1)
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		res, err := placeholderService.GetPlaceholder(ctx, placeholderID.String())
		assert.EqualError(t, err, "error")
		assert.Empty(t, res)
	})

	t.Run("negative - get single placeholder return error", func(t *testing.T) {
		mockPlaceholderCache.EXPECT().GetPlaceholderInfo(ctx, placeholderID.String()).Return(&model.PlaceholderDTO{}, redis.Nil).Times(1)
		mockPlaceholderRepo.EXPECT().GetSinglePlaceholder(ctx, placeholderID.String()).Return(model.PlaceholderDAO{}, errors.New("error")).Times(1)
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		res, err := placeholderService.GetPlaceholder(ctx, placeholderID.String())
		assert.EqualError(t, err, "error")
		assert.Empty(t, res)
	})

	t.Run("negative - placeholder not found", func(t *testing.T) {
		mockPlaceholderCache.EXPECT().GetPlaceholderInfo(ctx, placeholderID.String()).Return(&model.PlaceholderDTO{}, redis.Nil).Times(1)
		mockPlaceholderRepo.EXPECT().GetSinglePlaceholder(ctx, placeholderID.String()).Return(model.PlaceholderDAO{}, sql.ErrNoRows).Times(1)
		mockLogger.EXPECT().Info(gomock.Any()).Times(1)

		res, err := placeholderService.GetPlaceholder(ctx, placeholderID.String())
		assert.EqualError(t, err, common.ErrPlaceholderNotFound.Error())
		assert.Empty(t, res)
	})

	t.Run("negative - set placeholder to cache returning error", func(t *testing.T) {
		mockPlaceholderCache.EXPECT().GetPlaceholderInfo(ctx, placeholderID.String()).Return(&model.PlaceholderDTO{}, redis.Nil).Times(1)
		mockPlaceholderRepo.EXPECT().GetSinglePlaceholder(ctx, placeholderID.String()).Return(model.PlaceholderDAO{}, nil).Times(1)
		mockPlaceholderCache.EXPECT().SetPlaceholderInfo(ctx, gomock.Any()).Return(errors.New("error")).Times(1)
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		res, err := placeholderService.GetPlaceholder(ctx, placeholderID.String())
		assert.EqualError(t, err, "error")
		assert.Empty(t, res)
	})

	t.Run("positive - set cache when placeholder not found", func(t *testing.T) {
		mockPlaceholderCache.EXPECT().GetPlaceholderInfo(ctx, placeholderID.String()).Return(&model.PlaceholderDTO{}, redis.Nil).Times(1)
		mockPlaceholderRepo.EXPECT().GetSinglePlaceholder(ctx, placeholderID.String()).Return(model.PlaceholderDAO{}, nil).Times(1)
		mockPlaceholderCache.EXPECT().SetPlaceholderInfo(ctx, gomock.Any()).Return(nil).Times(1)
		mockAlphaProxy.EXPECT().GetPlaceholderStatus(ctx, gomock.Any()).Return(alpha.AlphaResponse{}, nil).Times(1)

		res, err := placeholderService.GetPlaceholder(ctx, placeholderID.String())
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})

}
