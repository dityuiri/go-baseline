package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/dityuiri/go-baseline/adapter/client/request"
)

const (
	defaultMaxRetry = 5

	// RequestIDHeaderName is a key used to get `X-Request-Id` from request header
	RequestIDHeaderName = "X-Request-Id"

	spanTitle = "Mapan_HTTP_Client"
)

func NewClient(ctx context.Context, config *Configuration, opts ...Option) IClient {
	id, _ := uuid.NewRandom()
	doer := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	client := &Client{
		Context: ctx,
		Config:  config,

		id:       id.String(),
		doer:     doer,
		maxRetry: defaultMaxRetry,
	}

	return client.set(opts...)
}

func (c *Client) set(opts ...Option) *Client {
	for _, fn := range opts {
		fn(c)
	}
	return c
}

func (c *Client) Close() {
	// todo : force close all connection
}

// Get will create new request object, calling Do to calls the server
func (c *Client) Get(url string, opts ...request.Option) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, url, nil, opts...)
}

// Post will create new request object, calling Do to calls the server
func (c *Client) Post(url string, body io.Reader, opts ...request.Option) (*http.Response, error) {
	return c.makeRequest(http.MethodPost, url, body, opts...)
}

// Put will create new request object, calling Do to calls the server
func (c *Client) Put(url string, body io.Reader, opts ...request.Option) (*http.Response, error) {
	return c.makeRequest(http.MethodPut, url, body, opts...)
}

// Patch will create new request object, calling Do to calls the server
func (c *Client) Patch(url string, body io.Reader, opts ...request.Option) (*http.Response, error) {
	return c.makeRequest(http.MethodPatch, url, body, opts...)
}

// Delete will create new request object, calling Do to calls the server
func (c *Client) Delete(url string, body io.Reader, opts ...request.Option) (*http.Response, error) {
	return c.makeRequest(http.MethodDelete, url, body, opts...)
}

func (c *Client) makeRequest(method string, url string, body io.Reader, opts ...request.Option) (*http.Response, error) {
	options := request.NewOptions(c.Context, opts...)

	if request, err := http.NewRequestWithContext(options.Context, method, url, body); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s request", method))
	} else {
		setRequestID := false

		request.Header = options.Headers

		requestID := options.RequestID

		if requestID == "" {
			requestID = request.Header.Get(RequestIDHeaderName)

			if requestID == "" {
				requestID = uuid.New().String()

				setRequestID = true
			}
		} else {
			setRequestID = true
		}

		if setRequestID {
			// This can cause the panic for concurrent map writes
			// Especially when using the same client in Go routines
			request.Header.Set(RequestIDHeaderName, requestID)
		}

		return c.Do(request)
	}
}

/*Do will call server with request
 * This functon will add request id in header if not exist
 */
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	var bodyReader *bytes.Reader

	if request.Body != nil {
		reqData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		request.Body = ioutil.NopCloser(bodyReader) // prevents closing the body between retries
	}

	multiErr := []error{}
	var response *http.Response

	for i := 0; i <= c.maxRetry; i++ {
		if response != nil {
			response.Body.Close()
		}

		response, _ = c.doer.Do(request)
		if bodyReader != nil {
			_, _ = bodyReader.Seek(0, 0)
		}

		multiErr = []error{}
		break
	}

	var errString string
	for i, e := range multiErr {
		errString += fmt.Sprintf("error in iteration %v: %v; ", i, e.Error())
	}

	var finalErr error
	if errString != "" {
		finalErr = fmt.Errorf(errString)
	}

	return response, finalErr
}

// Parse will parse the body from the reader and supports GZip compression.
// It will try to unmarshal the body from JSON.
func (c *Client) Parse(response *http.Response, data interface{}) error {
	defer response.Body.Close()

	body := response.Body

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		var err error

		body, err = gzip.NewReader(body)
		if err != nil {
			return err
		}
	}

	// This will strip away any additional information like encoding.
	mediatype, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))

	if body, err := ioutil.ReadAll(body); err != nil {
		return err
	} else if mediatype == "application/json" {
		return json.Unmarshal(body, data)
	} else {
		return errors.New(string(body))
	}
}
