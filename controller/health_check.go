package controller

import (
	"net/http"

	"github.com/dityuiri/go-adapter/logger"
	"github.com/dityuiri/go-baseline/common/util"
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

// Ping is a controller function to health check the app.
func (c *HealthCheckController) Ping(w http.ResponseWriter, r *http.Request) {
	util.WriteResponse(w, c.HealthCheckService.Ping(), http.StatusOK)
}
