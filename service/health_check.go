package service

//go:generate mockgen -package=service_mock -destination=../mock/service/health_check.go . IHealthCheckService

import (
	"github.com/dityuiri/go-baseline/repository"
)

type (
	// IHealthCheckService is an interface that has all the function to be implemented inside health check service
	IHealthCheckService interface {
		Ping() map[string]string
	}

	// HealthCheckService is a struct that consists of all the dependencies needed for health check service
	HealthCheckService struct {
		HealthCheckRepo repository.IHealthCheckRepository
	}
)

// Ping is a service function to do health check
func (s *HealthCheckService) Ping() map[string]string {
	var status = "OK"

	if err := s.HealthCheckRepo.Ping(); err != nil {
		status = "NOT OK"
	}

	return map[string]string{"status": status}
}
