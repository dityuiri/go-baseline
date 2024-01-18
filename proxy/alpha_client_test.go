package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dityuiri/go-baseline/adapter/client"
	clientMock "github.com/dityuiri/go-baseline/adapter/client/mock"
	loggerMock "github.com/dityuiri/go-baseline/adapter/logger/mock"
	"github.com/dityuiri/go-baseline/common"
	"github.com/dityuiri/go-baseline/config"
	"github.com/dityuiri/go-baseline/model/alpha"
)

func TestAlphaProxy_GetPlaceholderStatus(t *testing.T) {
	var (
		mockCtrl       = gomock.NewController(t)
		mockLogger     = loggerMock.NewMockILogger(mockCtrl)
		mockHTTPClient = clientMock.NewMockIClient(mockCtrl)

		proxy = AlphaProxy{
			Logger:     mockLogger,
			HTTPClient: mockHTTPClient,
			ClientConfiguration: config.HttpClient{
				ClientConfig: &client.Configuration{
					Timeout: 10,
				},
				ProxyURLs: config.ProxyURLs{
					AlphaURL: "localhost:8080",
				},
			},
		}

		ctx      = context.Background()
		alphaReq = alpha.AlphaRequest{
			PlaceholderID: uuid.New().String(),
			Amount:        25000,
		}

		output = alpha.AlphaResponse{
			ID:            uuid.New().String(),
			PlaceholderID: alphaReq.PlaceholderID,
			Amount:        alphaReq.Amount,
			Status:        "OK",
		}

		marshalledOutput, _ = json.Marshal(output)
		finalEndpoint       = fmt.Sprintf("%s%s", proxy.ClientConfiguration.ProxyURLs.AlphaURL, getPlaceholderStatus)
		header              = http.Header{}
	)

	t.Run("positive - statusOK", func(t *testing.T) {
		var (
			reader   = io.NopCloser(bytes.NewReader(marshalledOutput))
			response = &http.Response{
				Header:     header,
				Body:       reader,
				StatusCode: http.StatusOK,
			}
		)

		mockHTTPClient.EXPECT().Post(finalEndpoint, gomock.Any()).Return(response, nil).Times(1)

		header.Set("Content-Type", "application/json")
		res, err := proxy.GetPlaceholderStatus(ctx, alphaReq)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("parsing failed", func(t *testing.T) {
		var (
			reader   = io.NopCloser(strings.NewReader("potato"))
			response = &http.Response{
				Header:     header,
				Body:       reader,
				StatusCode: http.StatusOK,
			}
		)

		mockHTTPClient.EXPECT().Post(finalEndpoint, gomock.Any()).Return(response, nil).Times(1)
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		header.Set("Content-Type", "application/json")
		res, err := proxy.GetPlaceholderStatus(ctx, alphaReq)
		assert.NotNil(t, err)
		assert.Empty(t, res)
	})

	t.Run("status not found", func(t *testing.T) {
		var (
			reader   = io.NopCloser(bytes.NewReader(marshalledOutput))
			response = &http.Response{
				Header:     header,
				Body:       reader,
				StatusCode: http.StatusNotFound,
			}
		)

		mockHTTPClient.EXPECT().Post(finalEndpoint, gomock.Any()).Return(response, nil).Times(1)

		header.Set("Content-Type", "application/json")
		res, err := proxy.GetPlaceholderStatus(ctx, alphaReq)
		assert.EqualError(t, common.ErrAlphaProxyNotFound, err.Error())
		assert.NotEmpty(t, res)
	})

	t.Run("status bad request", func(t *testing.T) {
		var (
			reader   = io.NopCloser(bytes.NewReader(marshalledOutput))
			response = &http.Response{
				Header:     header,
				Body:       reader,
				StatusCode: http.StatusBadRequest,
			}
		)

		mockHTTPClient.EXPECT().Post(finalEndpoint, gomock.Any()).Return(response, nil).Times(1)

		header.Set("Content-Type", "application/json")
		res, err := proxy.GetPlaceholderStatus(ctx, alphaReq)
		assert.EqualError(t, common.ErrAlphaProxyBadRequest, err.Error())
		assert.NotEmpty(t, res)
	})

	t.Run("status internal server error", func(t *testing.T) {
		var (
			reader   = io.NopCloser(bytes.NewReader(marshalledOutput))
			response = &http.Response{
				Header:     header,
				Body:       reader,
				StatusCode: http.StatusInternalServerError,
			}
		)

		mockHTTPClient.EXPECT().Post(finalEndpoint, gomock.Any()).Return(response, nil).Times(1)

		header.Set("Content-Type", "application/json")
		res, err := proxy.GetPlaceholderStatus(ctx, alphaReq)
		assert.EqualError(t, common.ErrAlphaInternalServerError, err.Error())
		assert.NotEmpty(t, res)
	})

	t.Run("client post method error", func(t *testing.T) {
		mockHTTPClient.EXPECT().Post(finalEndpoint, gomock.Any()).Return(&http.Response{}, errors.New("error")).Times(1)
		mockLogger.EXPECT().Error(gomock.Any()).Times(1)

		header.Set("Content-Type", "application/json")
		res, err := proxy.GetPlaceholderStatus(ctx, alphaReq)
		assert.EqualError(t, err, "error")
		assert.Empty(t, res)
	})

	t.Run("json marshal error", func(t *testing.T) {
		jsonMarshal := json.Marshal
		common.JsonMarshal = func(v any) ([]byte, error) {
			return []byte{}, errors.New("error")
		}

		mockLogger.EXPECT().Error(gomock.Any()).Times(1)
		res, err := proxy.GetPlaceholderStatus(ctx, alphaReq)
		assert.EqualError(t, err, "error")
		assert.Empty(t, res)

		common.JsonMarshal = jsonMarshal
	})

}
