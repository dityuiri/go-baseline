package model

import "fmt"

type (
	APIResponse struct {
		Result interface{}       `json:"result,omitempty"`
		Error  *APIErrorResponse `json:"error,omitempty"`
	}

	APIErrorResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	APIErrorCode int
)

func (e APIErrorCode) String() string {
	return fmt.Sprintf("PLA-%03d", e)
}

const (
	// common
	InternalServerError APIErrorCode = iota + 001
	InvalidParameter
	MissingParameter
	InvalidRequestBody
	ObjectNotFound
)
