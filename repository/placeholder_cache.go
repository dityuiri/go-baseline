package repository

//go:generate mockgen -package=repository_mock -destination=../mock/repository/placeholder_cache.go . IPlaceholderCache

import (
    "context"
    "fmt"

    "github.com/dityuiri/go-baseline/adapter/logger"
    "github.com/dityuiri/go-baseline/adapter/redis"
    "github.com/dityuiri/go-baseline/model/dto"
)

type (
    IPlaceholderCache interface {
        SetAlphaPlaceholderInfo(ctx context.Context, placeholderDTO dto.PlaceholderDTO) error
        GetAlphaPlaceholderInfo(ctx context.Context, placeholderID string) (*dto.PlaceholderDTO, error)
    }

    PlaceholderCache struct {
        Redis redis.IRedis
        Logger logger.ILogger
    }
)

const (
    keyPlaceholder = "placeholder:%s"
)

func (pc *PlaceholderCache)  SetPlaceholderInfo(ctx context.Context, placeholderDTO dto.PlaceholderDTO) error {
    var key = fmt.Sprintf(keyPlaceholder, placeholderDTO.ID.String())

    err := pc.Redis.SetAsBytes(key, alphaResp)
    return err


func (pc *PlaceholderCache) GetPlaceholderInfo(ctx context.Context, placeholderID string) (*dto.PlaceholderDTO, error) {
    var (
        key = fmt.Sprintf(keyPlaceholder, placeholderID)
        result = &dto.PlaceholderDTO{}
    )

    err := apc.Redis.GetAndParseBytes(key, result)
    return result, err
}