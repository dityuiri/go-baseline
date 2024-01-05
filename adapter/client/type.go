package client

//go:generate mockgen -destination=mock/client.go -package=mock . IClient,Doer

import (
	"context"
	"io"
	"net/http"

	"github.com/dityuiri/go-baseline/adapter/client/request"
)

/*
Client implements IClient
stolen from https://github.com/gojek/heimdall
*/
type Client struct {
	Context context.Context
	Config  *Configuration

	id       string
	doer     Doer
	maxRetry int
}

// IClient wraps Doer, adding retries and circuit breaker
type IClient interface {
	Close()

	Get(url string, opts ...request.Option) (*http.Response, error)
	Post(url string, body io.Reader, opts ...request.Option) (*http.Response, error)
	Put(url string, body io.Reader, opts ...request.Option) (*http.Response, error)
	Patch(url string, body io.Reader, opts ...request.Option) (*http.Response, error)
	Delete(url string, body io.Reader, opts ...request.Option) (*http.Response, error)

	Do(request *http.Request) (*http.Response, error)

	Parse(response *http.Response, data interface{}) error
}

// Doer is the actual client that will "Do" the request vie "Do" function
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}
