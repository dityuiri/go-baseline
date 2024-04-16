package repository

import (
	"context"
	"github.com/dityuiri/go-baseline/model"

	"github.com/dityuiri/go-adapter/db"
	"github.com/dityuiri/go-adapter/logger"
)

//go:generate mockgen -package=repository_mock -destination=../mock/repository/placeholder_db.go . IPlaceholderRepository

type (
	IPlaceholderRepository interface {
		GetSinglePlaceholder(ctx context.Context, placeholderID string) (model.PlaceholderDAO, error)
		InsertPlaceholder(ctx context.Context, tx db.ITransaction, placeholder model.PlaceholderDAO) error
		UpdatePlaceholder(ctx context.Context, tx db.ITransaction, placeholder model.PlaceholderDAO) error
	}

	PlaceholderRepository struct {
		Logger logger.ILogger
		DB     db.IDatabase
	}
)

func (pr *PlaceholderRepository) GetSinglePlaceholder(ctx context.Context, placeholderID string) (model.PlaceholderDAO, error) {
	// Your implementation goes here
	return model.PlaceholderDAO{}, nil
}

func (pr *PlaceholderRepository) InsertPlaceholder(ctx context.Context, tx db.ITransaction, placeholder model.PlaceholderDAO) error {
	// Your implementation goes here
	return nil
}

func (pr *PlaceholderRepository) UpdatePlaceholder(ctx context.Context, tx db.ITransaction, placeholder model.PlaceholderDAO) error {
	// Your implementation goes here
	return nil
}
