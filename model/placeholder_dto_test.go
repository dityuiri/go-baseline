package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlaceholderCreateRequest_ToPlaceholderDTO(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var (
			input = PlaceholderCreateRequest{
				Name:   "Aoi",
				Amount: 10000,
			}

			expectedName   = input.Name
			expectedAmount = input.Amount
		)

		res := input.ToPlaceholderDTO()
		assert.NotNil(t, res.ID)
		assert.Equal(t, expectedName, res.Name)
		assert.Equal(t, expectedAmount, res.Amount)
	})
}

func TestPlaceholderDTO_ToPlaceholderDAO(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var (
			input = PlaceholderDTO{
				ID:        uuid.New(),
				Name:      "Minase",
				Amount:    25000,
				CreatedBy: "System",
				UpdatedBy: "System",
			}

			expected = PlaceholderDAO{
				ID:        input.ID,
				Name:      input.Name,
				Amount:    input.Amount,
				CreatedAt: "0001-01-01T00:00:00Z",
				UpdatedAt: "0001-01-01T00:00:00Z",
				CreatedBy: input.CreatedBy,
				UpdatedBy: input.UpdatedBy,
			}
		)

		res := input.ToPlaceholderDAO()
		assert.Equal(t, expected, res)
	})
}

func TestPlaceholderDTO_ToPlaceholderCreateResponse(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		var (
			input = PlaceholderDTO{
				ID:        uuid.New(),
				Name:      "Minase",
				Amount:    25000,
				CreatedBy: "System",
				UpdatedBy: "System",
			}

			expected = PlaceholderCreateResponse{
				ID:     input.ID.String(),
				Name:   input.Name,
				Amount: input.Amount,
			}
		)

		res := input.ToPlaceholderCreateResponse()
		assert.Equal(t, expected, res)
	})
}

func TestPlaceholderDTO_ToPlaceholderGetResponse(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		var (
			input = PlaceholderDTO{
				ID:        uuid.New(),
				Name:      "Minase",
				Amount:    25000,
				CreatedBy: "System",
				UpdatedBy: "System",
			}

			expected = PlaceholderGetResponse{
				ID:     input.ID.String(),
				Name:   input.Name,
				Amount: input.Amount,
			}
		)

		res := input.ToPlaceholderGetResponse()
		assert.Equal(t, expected, res)
	})
}
