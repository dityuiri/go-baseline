package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/dityuiri/go-baseline/model/dao"
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

	PlaceholderDTO struct {
		ID        uuid.UUID
		Name      string
		Amount    int
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

func (pDTO *PlaceholderDTO) ToPlaceholderDAO() dao.PlaceholderDAO {
	var (
		pDAO dao.PlaceholderDAO

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

//func (pDAO *PlaceholderDAO) ToPlaceholderDTO() dto.PlaceholderDTO {
//	var (
//		pDTO dto.PlaceholderDTO
//
//		a, _ = json.Marshal(pDAO)
//		_    = json.Unmarshal(a, &pDTO)
//	)
//
//	return pDTO
//}
