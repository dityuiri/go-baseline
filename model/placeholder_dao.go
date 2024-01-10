package model

import (
	"encoding/json"
	"github.com/google/uuid"
)

type (
	PlaceholderDAO struct {
		ID        uuid.UUID
		Name      string
		Amount    int
		CreatedAt interface{}
		CreatedBy string
		UpdatedAt interface{}
		UpdatedBy string
	}
)

func (pDAO *PlaceholderDAO) ToPlaceholderDTO() PlaceholderDTO {
	var (
		pDTO PlaceholderDTO

		b, _ = json.Marshal(pDAO)
		_    = json.Unmarshal(b, &pDTO)
	)

	return pDTO
}
