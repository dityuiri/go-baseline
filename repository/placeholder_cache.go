package repository

//go:generate mockgen -package=repository_mock -destination=../mock/repository/placeholder_cache.go . IPlaceholderCache

import (
	"context"
	"fmt"
	"github.com/dityuiri/go-baseline/model"

	"github.com/dityuiri/go-adapter/logger"
	"github.com/dityuiri/go-adapter/redis"
)

type (
	IPlaceholderCache interface {
		SetPlaceholderInfo(ctx context.Context, placeholderDTO model.PlaceholderDTO) error
		GetPlaceholderInfo(ctx context.Context, placeholderID string) (*model.PlaceholderDTO, error)
	}

	PlaceholderCache struct {
		Redis  redis.IRedis
		Logger logger.ILogger
	}
)

const (
	keyPlaceholder = "placeholder:%s"
)

func (pc *PlaceholderCache) SetPlaceholderInfo(ctx context.Context, placeholderDTO model.PlaceholderDTO) error {
	var key = fmt.Sprintf(keyPlaceholder, placeholderDTO.ID.String())

	err := pc.Redis.SetAsBytes(key, placeholderDTO)
	return err
}

func (pc *PlaceholderCache) GetPlaceholderInfo(ctx context.Context, placeholderID string) (*model.PlaceholderDTO, error) {
	var (
		key    = fmt.Sprintf(keyPlaceholder, placeholderID)
		result = &model.PlaceholderDTO{}
	)

	err := pc.Redis.GetAndParseBytes(key, result)
	return result, err
}
