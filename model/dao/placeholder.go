package dao

import (
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
