package repository

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	databaseMock "github.com/dityuiri/go-adapter/db/mock"
	loggerMock "github.com/dityuiri/go-adapter/logger/mock"
	"github.com/dityuiri/go-baseline/model"
)

func TestPlaceholderRepository_GetSinglePlaceholder(t *testing.T) {
	var (
		mockCtrl   = gomock.NewController(t)
		mockLogger = loggerMock.NewMockILogger(mockCtrl)
		mockDB     = databaseMock.NewMockIDatabase(mockCtrl)

		repo = PlaceholderRepository{
			Logger: mockLogger,
			DB:     mockDB,
		}

		ctx           = context.Background()
		placeholderID = uuid.New().String()
	)

	defer mockCtrl.Finish()

	t.Run("positive", func(t *testing.T) {
		res, err := repo.GetSinglePlaceholder(ctx, placeholderID)
		assert.Nil(t, err)
		assert.Empty(t, res)
	})
}

func TestPlaceholderRepository_InsertPlaceholder(t *testing.T) {
	var (
		mockCtrl   = gomock.NewController(t)
		mockLogger = loggerMock.NewMockILogger(mockCtrl)
		mockDB     = databaseMock.NewMockIDatabase(mockCtrl)
		mockTx     = databaseMock.NewMockITransaction(mockCtrl)

		repo = PlaceholderRepository{
			Logger: mockLogger,
			DB:     mockDB,
		}

		ctx         = context.Background()
		placeholder = model.PlaceholderDAO{
			ID:   uuid.New(),
			Name: "you know, a placeholder",
		}
	)

	defer mockCtrl.Finish()

	t.Run("positive", func(t *testing.T) {
		err := repo.InsertPlaceholder(ctx, mockTx, placeholder)
		assert.Nil(t, err)
	})
}

func TestPlaceholderRepository_UpdatePlaceholder(t *testing.T) {
	var (
		mockCtrl   = gomock.NewController(t)
		mockLogger = loggerMock.NewMockILogger(mockCtrl)
		mockDB     = databaseMock.NewMockIDatabase(mockCtrl)
		mockTx     = databaseMock.NewMockITransaction(mockCtrl)

		repo = PlaceholderRepository{
			Logger: mockLogger,
			DB:     mockDB,
		}

		ctx         = context.Background()
		placeholder = model.PlaceholderDAO{
			ID:   uuid.New(),
			Name: "you know, a placeholder",
		}
	)

	defer mockCtrl.Finish()

	t.Run("positive", func(t *testing.T) {
		err := repo.UpdatePlaceholder(ctx, mockTx, placeholder)
		assert.Nil(t, err)
	})
}
