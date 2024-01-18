package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlaceholderDAO_ToPlaceholderDTO(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var (
			input = PlaceholderDAO{
				ID:        uuid.New(),
				Name:      "Minase",
				Amount:    25000,
				CreatedBy: "System",
				UpdatedBy: "System",
			}

			expected = PlaceholderDTO{
				ID:        input.ID,
				Name:      input.Name,
				Amount:    input.Amount,
				CreatedBy: input.CreatedBy,
				UpdatedBy: input.UpdatedBy,
			}
		)

		res := input.ToPlaceholderDTO()
		assert.Equal(t, expected, res)
	})
}
