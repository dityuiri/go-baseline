package request

import (
	"context"
	"net/http"
)

type Option func(*Options)

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func WithHeaders(headers http.Header) Option {
	return func(o *Options) {
		o.Headers = headers
	}
}

func WithRequestID(requestID string) Option {
	return func(o *Options) {
		o.RequestID = requestID
	}
}
