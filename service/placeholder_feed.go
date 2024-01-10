package service

import (
	"context"

	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/common"
	"github.com/dityuiri/go-baseline/model"
	"github.com/dityuiri/go-baseline/repository"
)

//go:generate mockgen -package=service_mock -destination=../mock/service/placeholder_feed.go . IPlaceholderFeedService

type (
	IPlaceholderFeedService interface {
		PlaceholderRecorded(ctx context.Context, placeholderMsg model.PlaceholderMessage) (bool, error)
	}

	PlaceholderFeedService struct {
		Logger              logger.ILogger
		PlaceholderProducer repository.IPlaceholderProducer
	}
)

func (fs *PlaceholderFeedService) PlaceholderRecorded(ctx context.Context, placeholderMsg model.PlaceholderMessage) (bool, error) {
	placeholderMsg.EventName = common.EventPlaceholderRecorded
	err := fs.PlaceholderProducer.ProducePlaceholderRecord(ctx, placeholderMsg)
	if err != nil {
		fs.Logger.Error("failed to produce placeholder message")
		return false, err
	}

	return true, nil
}
