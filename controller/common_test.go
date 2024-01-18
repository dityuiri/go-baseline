package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dityuiri/go-baseline/common"
	"github.com/dityuiri/go-baseline/model"
)

func TestCommon_NewError(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		var (
			code = model.MissingParameter
			err  = common.ErrMissingPlaceholderID

			expected = model.APIResponse{
				Error: &model.APIErrorResponse{
					Code:    code.String(),
					Message: err.Error(),
				},
			}
		)

		result := NewError(code, err)
		assert.Equal(t, expected, result)
	})
}
