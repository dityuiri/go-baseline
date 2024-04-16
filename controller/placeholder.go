package controller

import (
	"github.com/dityuiri/go-baseline/model"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/dityuiri/go-adapter/logger"
	"github.com/dityuiri/go-baseline/common"
	"github.com/dityuiri/go-baseline/common/util"
	"github.com/dityuiri/go-baseline/service"
)

type (
	IPlaceholderController interface {
		GetPlaceholder(w http.ResponseWriter, r *http.Request)
		CreatePlaceholder(w http.ResponseWriter, r *http.Request)
	}

	PlaceholderController struct {
		Logger             logger.ILogger
		PlaceholderService service.IPlaceholderService
	}
)

func (c *PlaceholderController) GetPlaceholder(w http.ResponseWriter, r *http.Request) {
	var (
		resp = model.APIResponse{}
		ctx  = r.Context()

		queryPlaceholderID = r.URL.Query().Get("placeholder_id")
	)

	if queryPlaceholderID == "" {
		errResponse := NewError(model.MissingParameter, common.ErrMissingPlaceholderID)
		util.WriteResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(strings.TrimSpace(queryPlaceholderID))
	if err != nil {
		errResponse := NewError(model.InvalidParameter, common.ErrInvalidUUIDPlaceholderID)
		util.WriteResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	result, err := c.PlaceholderService.GetPlaceholder(ctx, queryPlaceholderID)
	if err != nil {
		var (
			status = http.StatusInternalServerError
			code   = model.InternalServerError
		)

		switch err {
		case common.ErrPlaceholderNotFound, common.ErrAlphaProxyNotFound:
			status = http.StatusNotFound
			code = model.ObjectNotFound
		}

		errResponse := NewError(code, err)
		util.WriteResponse(w, errResponse, status)
		return
	}

	resp.Result = map[string]interface{}{
		common.PlaceholderKey: result,
	}

	util.WriteResponse(w, resp, http.StatusOK)
}

func (c *PlaceholderController) CreatePlaceholder(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		resp = model.APIResponse{}

		placeholderCreateRequest *model.PlaceholderCreateRequest
	)

	if err := util.HttpRequestBodyParser(r, &placeholderCreateRequest); err != nil {
		errResponse := NewError(model.InvalidRequestBody, common.ErrInvalidRequestBody)
		util.WriteResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	createResponse, err := c.PlaceholderService.CreateNewPlaceholder(ctx, *placeholderCreateRequest)
	if err != nil {
		var (
			status = http.StatusInternalServerError
			code   = model.InternalServerError
		)

		errResponse := NewError(code, err)
		util.WriteResponse(w, errResponse, status)
		return
	}

	resp.Result = map[string]interface{}{
		common.PlaceholderKey: createResponse,
	}

	util.WriteResponse(w, resp, http.StatusOK)

}
