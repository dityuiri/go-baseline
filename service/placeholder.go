package service

import (
	"context"

	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/model/dto"
	"github.com/dityuiri/go-baseline/model/dto/alpha"
	"github.com/dityuiri/go-baseline/repository"
)

//go:generate mockgen -package=service_mock -destination=../mock/service/placeholder.go . IPlaceholderService

type (
	IPlaceholderService interface {
		CreateNewPlaceholder(ctx context.Context, placeholderRequest dto.PlaceholderCreateRequest) (dto.PlaceholderCreateResponse, error)
		GetPlaceholder(ctx context.Context, placeholderID string) (dto.PlaceholderGetResponse, error)
	}

	PlaceholderService struct {
		Logger                logger.ILogger
		PlaceholderRepository repository.IPlaceholderRepository
		PlaceholderCache      repository.IPlaceholderCache
		AlphaProxy            proxy.IAlphaProxy
	}
)

func (ps *PlaceholderService) CreateNewPlaceholder(ctx context.Context, placeholderRequest dto.PlaceholderCreateRequest) (dto.PlaceholderCreateResponse, error) {
	// Implement your code here
	var response dto.PlaceholderCreateResponse
	// Insert placeholder
	placeholderDTO := placeholderRequest.ToPlaceholderDTO()
	err := ps.PlaceholderRepository.InsertPlaceholder(ctx, nil, placeholderDTO.ToPlaceholderDAO())
	if err != nil {
		ps.Logger.Error("error inserting placeholder")
		return response, err
	}

	return placeholderDTO.ToPlaceholderCreateResponse(), err
}

func (ps *PlaceholderService) GetPlaceholder(ctx context.Context, placeholderID string) (dto.PlaceholderGetResponse, error) {
    // Implement your code here
    // Redis cache + db repository + http proxy example
    var placeholderResp dto.PlaceholderGetResponse

    // Try to get from the redis first
    placeholderDTO, err := ps.PlaceholderCache.GetPlaceholderInfo(ctx, placeholderID)
    if err != nil {
        if err != redis.Nil{
            ps.Logger.Error("error getting placeholder cache from redis")
            return placeholderResp, err
        }

        // Case key not found in redis
        // Proceed to get from db
        placeholderDTO, err := ps.PlaceholderRepository.GetSinglePlaceholder(ctx, placeholderID)
        if err != nil {
            if sql.ErrNoRows {
                ps.Logger.Info(fmt.Sprintf("placeholder with id %s not found", placeholderID))
                err = config.ErrPlaceholderNotFound
            }

            ps.Logger.Error("error getting placeholder data from db")
            return placeholderResp, err
        }

        // Set to redis
        err = ps.SetPlaceholderInfo(ctx, placeholderDTO)
        if err != nil {
            ps.Logger.Error("error set placeholder to redis cache")
            return placeholderResp, err
        }
    }

    placeholderResp = placeholderDTO.ToPlaceholderGetResponse()

    // Call to alpha proxy to get the placeholder status
    alphaReq := mapPlaceholderDTOToAlphaStatusRequest(placeholderDTO, )
    alphaResp, err := ps.AlphaProxy.GetPlaceholderStatus(ctx, alphaReq)
    if err != nil {
        return placeholderResp, err
    }

    // Assign status from Alpha
    placeholderResp.Status = alphaResp.Status
    return placeholderResp, err
}

func (ps *PlaceholderService) mapPlaceholderDTOToAlphaStatusRequest(placeholderDTO dto.PlaceholderDTO) alphaReq alpha.AlphaRequest{
    return alpha.AlphaRequest{
        PlaceholderID: placeholderDTO.ID.String(),
        Amount: placeholderDTO.Amount,
    }
}
