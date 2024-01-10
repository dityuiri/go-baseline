package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	// PlaceholderCreateRequest POST /v1/placeholder request
	PlaceholderCreateRequest struct {
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	// PlaceholderCreateResponse POST /v1/placeholder response
	PlaceholderCreateResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Amount int    `json:"amount"`
	}

	// PlaceholderGetResponse GET /v1/placehodler response
	PlaceholderGetResponse struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Amount int    `json:"amount"`
		Status string `json:"status"`
	}

	// PlaceholderMessage Message to be exchanged via Kafka / message broker
	PlaceholderMessage struct {
		ID        string `json:"id"`
		EventName string `json:"event_name"`
		Name      string `json:"name"`
		Amount    int    `json:"amount"`
		CreatedBy string `json:"created_by"`
		UpdatedBy string `json:"updated_by"`
		Error     string `json:"error"`
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
