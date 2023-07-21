package application

import (
	"stockbit-challenge/adapter/kafka/producer"
	"stockbit-challenge/repository"
	"stockbit-challenge/service"
)

type Dependency struct {
	TransactionFeedService service.ITransactionFeedService
	StockService           service.IStockService
}

func SetupDependency(app *App) *Dependency {
	stockRepo := &repository.StockRepository{
		Redis: app.Redis,
	}

	trxProducer := &repository.TransactionProducer{
		Producer:    producer.NewProducer(app.Context, app.Config.Kafka.Producer),
		KafkaConfig: app.Config.Kafka,
	}

	trxFeedService := &service.TransactionFeedService{
		StockRepository:     stockRepo,
		TransactionProducer: trxProducer,
	}

	stockService := &service.StockService{
		StockRepository: stockRepo,
	}

	return &Dependency{
		TransactionFeedService: trxFeedService,
		StockService:           stockService,
	}
}
