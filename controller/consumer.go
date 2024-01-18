package controller

import (
	"context"
	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/common"

	"github.com/dityuiri/go-baseline/model"
	"github.com/dityuiri/go-baseline/service"
)

type ConsumerHandler struct {
	Logger                 logger.ILogger
	PlaceholderFeedService service.IPlaceholderFeedService
}

func (ch *ConsumerHandler) Placeholder(msg []byte) (bool, error) {
	var (
		ctx         = context.Background()
		placeholder = &model.PlaceholderMessage{}
	)

	err := common.JsonUnmarshal(msg, placeholder)
	if err != nil {
		ch.Logger.Error("error unmarshalling message")
		return true, err
	}

	return ch.PlaceholderFeedService.PlaceholderRecorded(ctx, *placeholder)
}
