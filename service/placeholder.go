package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"

	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/common"
	"github.com/dityuiri/go-baseline/model"
	"github.com/dityuiri/go-baseline/model/alpha"
	"github.com/dityuiri/go-baseline/proxy"
	"github.com/dityuiri/go-baseline/repository"
)

//go:generate mockgen -package=service_mock -destination=../mock/service/placeholder.go . IPlaceholderService

type (
	IPlaceholderService interface {
		CreateNewPlaceholder(ctx context.Context, placeholderRequest model.PlaceholderCreateRequest) (model.PlaceholderCreateResponse, error)
		GetPlaceholder(ctx context.Context, placeholderID string) (model.PlaceholderGetResponse, error)
	}

	PlaceholderService struct {
		Logger                logger.ILogger
		PlaceholderRepository repository.IPlaceholderRepository
		PlaceholderCache      repository.IPlaceholderCache
		AlphaProxy            proxy.IAlphaProxy
	}
)

func (ps *PlaceholderService) CreateNewPlaceholder(ctx context.Context, placeholderRequest model.PlaceholderCreateRequest) (model.PlaceholderCreateResponse, error) {
	// Implement your code here
	var response model.PlaceholderCreateResponse
	// Insert placeholder
	placeholderDTO := placeholderRequest.ToPlaceholderDTO()
	err := ps.PlaceholderRepository.InsertPlaceholder(ctx, nil, placeholderDTO.ToPlaceholderDAO())
	if err != nil {
		ps.Logger.Error("error inserting placeholder")
		return response, err
	}

	return placeholderDTO.ToPlaceholderCreateResponse(), err
}

func (ps *PlaceholderService) GetPlaceholder(ctx context.Context, placeholderID string) (model.PlaceholderGetResponse, error) {
	// Implement your code here
	// Redis cache + db repository + http proxy example
	var placeholderResp model.PlaceholderGetResponse

	// Try to get from the redis first
	placeholderDTO, err := ps.PlaceholderCache.GetPlaceholderInfo(ctx, placeholderID)
	if err != nil {
		if err != redis.Nil {
			ps.Logger.Error("error getting placeholder cache from redis")
			return placeholderResp, err
		}

		// Case key not found in redis
		// Proceed to get from db
		placeholderDAO, err := ps.PlaceholderRepository.GetSinglePlaceholder(ctx, placeholderID)
		if err != nil {
			if err == sql.ErrNoRows {
				ps.Logger.Info(fmt.Sprintf("placeholder with id %s not found", placeholderID))
				err = common.ErrPlaceholderNotFound
			} else {
				ps.Logger.Error("error getting placeholder data from db")
			}

			return placeholderResp, err
		}

		placeholderFromDB := placeholderDAO.ToPlaceholderDTO()
		placeholderDTO = &placeholderFromDB

		// Set to redis
		err = ps.PlaceholderCache.SetPlaceholderInfo(ctx, placeholderFromDB)
		if err != nil {
			ps.Logger.Error("error set placeholder to redis cache")
			return placeholderResp, err
		}
	}

	placeholderResp = placeholderDTO.ToPlaceholderGetResponse()

	// Call to alpha proxy to get the placeholder status
	alphaReq := ps.mapPlaceholderDTOToAlphaStatusRequest(*placeholderDTO)
	alphaResp, err := ps.AlphaProxy.GetPlaceholderStatus(ctx, alphaReq)
	if err != nil {
		ps.Logger.Error("get placeholder status error")
		return placeholderResp, err
	}

	// Assign status from Alpha
	placeholderResp.Status = alphaResp.Status
	return placeholderResp, err
}

func (ps *PlaceholderService) mapPlaceholderDTOToAlphaStatusRequest(placeholderDTO model.PlaceholderDTO) alpha.AlphaRequest {
	return alpha.AlphaRequest{
		PlaceholderID: placeholderDTO.ID.String(),
		Amount:        placeholderDTO.Amount,
	}
}
