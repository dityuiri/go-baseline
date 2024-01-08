package proxy

//go:generate mockgen -package=proxy_mock -destination=../mock/proxy/alpha.go . IAlphaProxy

import (
    "json"

    "github.com/dityuiri/go-baseline/adapter/logger"
    "github.com/dityuiri/go-baseline/adapter/client"
    "github.com/dityuiri/go-baseline/model/dto/alpha"
)

type (
    IAlphaProxy interface {
        GetPlaceholderStatus(ctx context.Context, alpha.AlphaRequest) (alpha.AlphaResponse, error)
    }

    AlphaProxy struct {
        Logger logger.ILogger
        HTTPClient client.IClient
        Config  *config.Configuration
    }
)

func (ap *AlphaProxy) GetPlaceholderStatus(ctx context.Context, alphaReq alpha.AlphaRequest) (alpha.AlphaResponse, error) {
    var (
        result = &alpha.AlphaResponse{}
        finalEndpoint = fmt.Sprintf("%s%s", ap.Config.Constants.Proxy.AlphaURL, ap.Config.Constants.Proxy.AlphaGetPlaceholderStatusEndpoint)
    )

    reqOut, err := json.Marshal(alphaReq)
    if err != nil {
        ap.Logger.Error("error marshaling request")
        return result, err
    }

    resp, err := ap.Post(finalEndpoint, bytes.NewBuffer(reqOut))
    if err != nil {
        ap.Logger.Error("error executing POST request to Alpha")
        return result, err
    }

    switch(resp.StatusCode) {
    case http.StatusOK:
        err = json.Unmarshal(resp, result)
        if err != nil {
            ap.Logger.Error("error unmarshalling response from alpha")
        }
        return result, err
    case http.StatusNotFound:
        return result, config.ErrAlphaProxyNotFound
    case http.StatusBadRequest:
        return result, config.ErrAlphaProxyBadRequest
    default:
        return result, err
    }

    return result, err
}