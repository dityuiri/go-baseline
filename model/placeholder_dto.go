package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	PlaceholderCreateRequest struct {
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	PlaceholderCreateResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	PlaceholderGetResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Amount int    `json:"amount"`
		Status string `json:"status"`
	}

	PlaceholderDTO struct {
		ID        uuid.UUID
		Name      string
		Amount    int
		Status    string
		CreatedAt time.Time
		CreatedBy string
		UpdatedAt time.Time
		UpdatedBy string
	}
)

func (pcr *PlaceholderCreateRequest) ToPlaceholderDTO() PlaceholderDTO {
	return PlaceholderDTO{
		ID:     uuid.New(),
		Name:   pcr.Name,
		Amount: pcr.Amount,
	}
}

func (pDTO *PlaceholderDTO) ToPlaceholderDAO() PlaceholderDAO {
	var (
		pDAO PlaceholderDAO

		b, _ = json.Marshal(pDTO)
		_    = json.Unmarshal(b, &pDAO)
	)

	return pDAO
}

func (pDTO *PlaceholderDTO) ToPlaceholderCreateResponse() PlaceholderCreateResponse {
	return PlaceholderCreateResponse{
		ID:     pDTO.ID.String(),
		Name:   pDTO.Name,
		Amount: pDTO.Amount,
	}
}

func (pDTO *PlaceholderDTO) ToPlaceholderGetResponse() PlaceholderGetResponse {
	return PlaceholderGetResponse{
		ID:     pDTO.ID.String(),
		Name:   pDTO.Name,
		Amount: pDTO.Amount,
	}
}
