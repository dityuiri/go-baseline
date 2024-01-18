package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIErrorCode_String(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		result := InternalServerError.String()
		assert.Equal(t, "PLA-001", result)
	})
}
