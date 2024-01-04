package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/service"
)

type (
	// IHealthCheckController is an interface ...
	IHealthCheckController interface {
		Ping(w http.ResponseWriter, r *http.Request)
	}

	// HealthCheckController is an app health check struct that consists of all the dependencies needed for health check controller
	HealthCheckController struct {
		Logger             logger.ILogger
		HealthCheckService service.IHealthCheckService
	}
)

func writeResponse(w http.ResponseWriter, body interface{}, status int) {
	resp, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}

// Ping is a controller function to health check the app.
func (c *HealthCheckController) Ping(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, c.HealthCheckService.Ping(), http.StatusOK)
}
