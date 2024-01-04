package controller

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"

	"github.com/dityuiri/go-baseline/mock"
	serviceMock "github.com/dityuiri/go-baseline/mock/service"
)

func TestHealthCheck_Ping(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var (
			mockHealthCheck = serviceMock.NewMockIHealthCheckService(mockCtrl)
			mockWriter      = mock.NewMockResponseWriter(mockCtrl)

			h = &HealthCheckController{
				HealthCheckService: mockHealthCheck,
			}
		)

		router := chi.NewRouter()
		router.Get("/ping", h.Ping)

		mockHealthCheck.EXPECT().Ping().Return(map[string]string{"status": "OK"})
		mockWriter.EXPECT().Header().Return(http.Header{})
		mockWriter.EXPECT().WriteHeader(http.StatusOK)
		mockWriter.EXPECT().Write(gomock.Any()).Return(0, nil)

		request, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(mockWriter, request)
	})
}
