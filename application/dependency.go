package application

import (
	"github.com/dityuiri/go-baseline/adapter/kafka/producer"
	"github.com/dityuiri/go-baseline/repository"
	"github.com/dityuiri/go-baseline/service"
)

type Dependency struct {
	TransactionFeedService service.ITransactionFeedService
	StockService           service.IStockService
	HealthCheckService     service.IHealthCheckService
	PlaceholderService     service.IPlaceholderService
}

func SetupDependency(app *App) *Dependency {
	stockRepo := &repository.StockRepository{
		Redis: app.Redis,
	}

	healthCheckRepo := &repository.HealthCheckRepository{}

	trxProducer := &repository.TransactionProducer{
		Producer:    producer.NewProducer(app.Config.Kafka.Producer),
		KafkaConfig: app.Config.Kafka,
	}

	placeholderRepo := &repository.PlaceholderRepository{
		Logger: app.Logger,
		DB:     app.DB,
	}

	trxFeedService := &service.TransactionFeedService{
		StockRepository:     stockRepo,
		TransactionProducer: trxProducer,
	}

	stockService := &service.StockService{
		StockRepository: stockRepo,
	}

	healthCheckService := &service.HealthCheckService{
		HealthCheckRepo: healthCheckRepo,
	}

	placeholderService := &service.PlaceholderService{
		Logger:                app.Logger,
		PlaceholderRepository: placeholderRepo,
	}

	return &Dependency{
		TransactionFeedService: trxFeedService,
		StockService:           stockService,
		HealthCheckService:     healthCheckService,
		PlaceholderService:     placeholderService,
	}
}
