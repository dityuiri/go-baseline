package controller

import (
	"github.com/dityuiri/go-baseline/model"
)

func NewError(code model.APIErrorCode, err error) model.APIResponse {
	return model.APIResponse{
		Error: &model.APIErrorResponse{
			Code:    code.String(),
			Message: err.Error(),
		},
	}
}
