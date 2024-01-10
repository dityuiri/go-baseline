package proxy

//go:generate mockgen -package=proxy_mock -destination=../mock/proxy/alpha.go . IAlphaProxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dityuiri/go-baseline/model/alpha"
	"net/http"

	"github.com/dityuiri/go-baseline/adapter/client"
	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/common"
	"github.com/dityuiri/go-baseline/common/util"
	"github.com/dityuiri/go-baseline/config"
)

type (
	IAlphaProxy interface {
		GetPlaceholderStatus(ctx context.Context, alphaReq alpha.AlphaRequest) (alpha.AlphaResponse, error)
	}

	AlphaProxy struct {
		Logger              logger.ILogger
		HTTPClient          client.IClient
		ClientConfiguration config.HttpClient
	}
)

const (
	getPlaceholderStatus = "/v1/placeholder/status"
)

func (ap *AlphaProxy) GetPlaceholderStatus(ctx context.Context, alphaReq alpha.AlphaRequest) (alpha.AlphaResponse, error) {
	var (
		result        = &alpha.AlphaResponse{}
		header        = http.Header{}
		finalEndpoint = fmt.Sprintf("%s%s", ap.ClientConfiguration.ProxyURLs.AlphaURL, getPlaceholderStatus)
	)

	reqOut, err := json.Marshal(alphaReq)
	if err != nil {
		ap.Logger.Error("error marshaling request")
		return *result, err
	}

	header.Set("Accept", "application/json, text/plain, */*")
	resp, err := ap.HTTPClient.Post(finalEndpoint, bytes.NewBuffer(reqOut))
	if err != nil {
		ap.Logger.Error("error executing POST request to Alpha")
		return *result, err
	}

	if err = util.HttpResponseBodyParser(resp, result); err != nil {
		return *result, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return *result, nil
	case http.StatusNotFound:
		return *result, common.ErrAlphaProxyNotFound
	case http.StatusBadRequest:
		return *result, common.ErrAlphaProxyBadRequest
	default:
		return *result, common.ErrAlphaInternalServerError
	}
}
