package service

import (
	"context"

	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/model/dto"
	"github.com/dityuiri/go-baseline/repository"
)

//go:generate mockgen -package=service_mock -destination=../mock/service/placeholder.go . IPlaceholderService

type (
	IPlaceholderService interface {
		CreateNewPlaceholder(ctx context.Context, placeholderRequest dto.PlaceholderCreateRequest) (dto.PlaceholderCreateResponse, error)
	}

	PlaceholderService struct {
		Logger                logger.ILogger
		PlaceholderRepository repository.IPlaceholderRepository
	}
)

func (ps *PlaceholderService) CreateNewPlaceholder(ctx context.Context, placeholderRequest dto.PlaceholderCreateRequest) (dto.PlaceholderCreateResponse, error) {
	// Implement your code here
	var response dto.PlaceholderCreateResponse
	// Insert placeholder
	placeholderDTO := placeholderRequest.ToPlaceholderDTO()
	err := ps.PlaceholderRepository.InsertPlaceholder(ctx, nil, placeholderDTO.ToPlaceholderDAO())
	if err != nil {
		ps.Logger.Error("Error")
		return response, err
	}

	return placeholderDTO.ToPlaceholderCreateResponse(), err
}
