package common

import "errors"

var (
	// Parameter Validation Errors
	ErrInvalidUUIDPlaceholderID = errors.New("invalid uuid #{placeholderID}")
	ErrInvalidRequestBody       = errors.New("invalid request body")
	ErrMissingPlaceholderID     = errors.New("missing #{placeholderID}")

	// Proxy Errors
	ErrAlphaProxyBadRequest     = errors.New("bad request from alpha")
	ErrAlphaProxyNotFound       = errors.New("alpha returned not found")
	ErrAlphaInternalServerError = errors.New("internal server error from alpha")

	// Repository Errors
	ErrPlaceholderNotFound = errors.New("placeholder not found")
)
