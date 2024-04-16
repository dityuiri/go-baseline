package controller

import (
	"errors"
	"fmt"
	"github.com/dityuiri/go-baseline/common"
	"github.com/dityuiri/go-baseline/model"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	loggerMock "github.com/dityuiri/go-adapter/logger/mock"
	"github.com/dityuiri/go-baseline/mock"
	serviceMock "github.com/dityuiri/go-baseline/mock/service"
)

func TestPlaceholderController_GetPlaceholder(t *testing.T) {
	var (
		mockCtrl               = gomock.NewController(t)
		mockLogger             = loggerMock.NewMockILogger(mockCtrl)
		mockWriter             = mock.NewMockResponseWriter(mockCtrl)
		mockPlaceholderService = serviceMock.NewMockIPlaceholderService(mockCtrl)

		placeholderController = PlaceholderController{
			Logger:             mockLogger,
			PlaceholderService: mockPlaceholderService,
		}

		baseRoute     = "/v1/placeholder"
		placeholderID = uuid.New()
		url           = fmt.Sprintf("%s?placeholder_id=%v", baseRoute, placeholderID)
		router        = chi.NewRouter()
	)

	defer mockCtrl.Finish()

	t.Run("positive", func(t *testing.T) {
		router.Get(baseRoute, placeholderController.GetPlaceholder)

		mockWriter.EXPECT().Header().Return(http.Header{})
		mockWriter.EXPECT().WriteHeader(http.StatusOK)
		mockWriter.EXPECT().Write(gomock.Any())
		mockPlaceholderService.EXPECT().GetPlaceholder(gomock.Any(), placeholderID.String()).Return(model.PlaceholderGetResponse{}, nil)

		request, _ := http.NewRequest("GET", url, nil)
		router.ServeHTTP(mockWriter, request)
	})

	t.Run("missing placeholder id", func(t *testing.T) {
		router.Get(baseRoute, placeholderController.GetPlaceholder)

		mockWriter.EXPECT().Header().Return(http.Header{})
		mockWriter.EXPECT().WriteHeader(http.StatusBadRequest)
		mockWriter.EXPECT().Write(gomock.Any())

		request, _ := http.NewRequest("GET", baseRoute, nil)
		router.ServeHTTP(mockWriter, request)
	})

	t.Run("non uuid placeholder id", func(t *testing.T) {
		invalidUrl := fmt.Sprintf("%s?placeholder_id=%v", baseRoute, "sausage")
		router.Get(baseRoute, placeholderController.GetPlaceholder)

		mockWriter.EXPECT().Header().Return(http.Header{})
		mockWriter.EXPECT().WriteHeader(http.StatusBadRequest)
		mockWriter.EXPECT().Write(gomock.Any())

		request, _ := http.NewRequest("GET", invalidUrl, nil)
		router.ServeHTTP(mockWriter, request)
	})

	t.Run("object not found", func(t *testing.T) {
		router.Get(baseRoute, placeholderController.GetPlaceholder)

		mockWriter.EXPECT().Header().Return(http.Header{})
		mockWriter.EXPECT().WriteHeader(http.StatusNotFound)
		mockWriter.EXPECT().Write(gomock.Any())
		mockPlaceholderService.EXPECT().GetPlaceholder(gomock.Any(), placeholderID.String()).Return(model.PlaceholderGetResponse{}, common.ErrPlaceholderNotFound)

		request, _ := http.NewRequest("GET", url, nil)
		router.ServeHTTP(mockWriter, request)
	})

	t.Run("internal server error", func(t *testing.T) {
		router.Get(baseRoute, placeholderController.GetPlaceholder)

		mockWriter.EXPECT().Header().Return(http.Header{})
		mockWriter.EXPECT().WriteHeader(http.StatusInternalServerError)
		mockWriter.EXPECT().Write(gomock.Any())
		mockPlaceholderService.EXPECT().GetPlaceholder(gomock.Any(), placeholderID.String()).Return(model.PlaceholderGetResponse{}, errors.New("error"))

		request, _ := http.NewRequest("GET", url, nil)
		router.ServeHTTP(mockWriter, request)
	})
}
