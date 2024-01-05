package request

import (
	"context"
	"net/http"
)

type Options struct {
	Context context.Context
	Headers http.Header

	RequestID string
}

func NewOptions(ctx context.Context, opts ...Option) *Options {
	options := &Options{
		Context: ctx,
		Headers: http.Header{},
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
