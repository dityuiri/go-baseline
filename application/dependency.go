package application

import (
	"github.com/dityuiri/go-adapter/client"
	"github.com/dityuiri/go-adapter/kafka/producer"
	"github.com/dityuiri/go-baseline/proxy"
	"github.com/dityuiri/go-baseline/repository"
	"github.com/dityuiri/go-baseline/service"
)

type Dependency struct {
	HealthCheckService     service.IHealthCheckService
	PlaceholderService     service.IPlaceholderService
	PlaceholderFeedService service.IPlaceholderFeedService
}

func SetupDependency(app *App) *Dependency {
	// Repository and Proxy layer
	healthCheckRepo := &repository.HealthCheckRepository{}
	//trxProducer := &repository.TransactionProducer{
	//	Producer:    producer.NewProducer(app.Config.Kafka.Producer),
	//	KafkaConfig: app.Config.Kafka,
	//}

	placeholderProducer := &repository.PlaceholderProducer{
		Producer:    producer.NewProducer(app.Config.Kafka.Producer),
		KafkaConfig: app.Config.Kafka,
	}

	placeholderRepo := &repository.PlaceholderRepository{
		Logger: app.Logger,
		DB:     app.DB,
	}

	placeholderCache := &repository.PlaceholderCache{
		Redis:  app.Redis,
		Logger: app.Logger,
	}

	alphaProxy := &proxy.AlphaProxy{
		Logger:     app.Logger,
		HTTPClient: client.NewClient(app.Context, app.Config.HTTPClient.ClientConfig),
	}

	// Service layer

	healthCheckService := &service.HealthCheckService{
		HealthCheckRepo: healthCheckRepo,
	}

	placeholderService := &service.PlaceholderService{
		Logger:                app.Logger,
		PlaceholderRepository: placeholderRepo,
		PlaceholderCache:      placeholderCache,
		AlphaProxy:            alphaProxy,
	}

	placeholderFeedService := &service.PlaceholderFeedService{
		Logger:              app.Logger,
		PlaceholderProducer: placeholderProducer,
	}

	return &Dependency{
		HealthCheckService:     healthCheckService,
		PlaceholderService:     placeholderService,
		PlaceholderFeedService: placeholderFeedService,
	}
}
