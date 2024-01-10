package controller

import (
	"context"
	"encoding/json"

	"github.com/dityuiri/go-baseline/model"
	"github.com/dityuiri/go-baseline/service"
)

type ConsumerHandler struct {
	PlaceholderFeedService service.IPlaceholderFeedService
}

func (ch *ConsumerHandler) Placeholder(msg []byte) (bool, error) {
	var (
		ctx         = context.Background()
		placeholder = &model.PlaceholderMessage{}
	)

	err := json.Unmarshal(msg, placeholder)
	if err != nil {
		return true, err
	}

	return ch.PlaceholderFeedService.PlaceholderRecorded(ctx, *placeholder)
}
