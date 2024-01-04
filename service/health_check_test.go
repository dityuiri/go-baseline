package service

//go:generate mockgen -package=service_mock -destination=../mock/service/health_check.go . IHealthCheckService

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	repositoryMock "github.com/dityuiri/go-baseline/mock/repository"
)

func TestPing(t *testing.T) {
	t.Run("ping success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepository := repositoryMock.NewMockIHealthCheckRepository(mockCtrl)

		mockRepository.EXPECT().Ping().Return(nil).Times(1)

		s := HealthCheckService{
			HealthCheckRepo: mockRepository,
		}

		result := s.Ping()
		assert.Equal(t, result["status"], "OK")
	})

	t.Run("ping error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepository := repositoryMock.NewMockIHealthCheckRepository(mockCtrl)

		mockRepository.EXPECT().Ping().Return(errors.New("error")).Times(1)

		s := HealthCheckService{
			HealthCheckRepo: mockRepository,
		}

		result := s.Ping()
		assert.Equal(t, result["status"], "NOT OK")
	})
}
