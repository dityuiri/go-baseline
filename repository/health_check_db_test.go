package repository

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Run("ping success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		r := HealthCheckRepository{}

		err := r.Ping()
		assert.Equal(t, err, nil)
	})
}
