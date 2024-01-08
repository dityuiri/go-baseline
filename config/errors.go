package config

import "errors"

const (
    ErrAlphaProxyBadRequest = errors.New("bad request from alpha")
    ErrAlphaProxyNotFound   = errors.New("alpha returned not found")
    ErrAlphaInternalServerError = errors.New("internal server error from alpha")
    ErrPlaceholderNotFound = errors.New("placeholder not found")
)