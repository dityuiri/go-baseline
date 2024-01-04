package repository

//go:generate mockgen -package=repository_mock -destination=../mock/repository/health_check.go . IHealthCheckRepository

type (
	// IHealthCheckRepository is a repository interface that consists of repository functions for app health checking
	IHealthCheckRepository interface {
		Ping() error
	}

	// HealthCheckRepository is a struct that consists of App that has Database adapter inside
	HealthCheckRepository struct {
	}
)

// Ping is a repository function for app health checking
func (r *HealthCheckRepository) Ping() error {
	return nil
}
