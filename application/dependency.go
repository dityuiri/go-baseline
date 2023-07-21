package application

import (
	"stockbit-challenge/adapter/kafka/producer"
	"stockbit-challenge/repository"
	"stockbit-challenge/service"
)

type Dependency struct {
	StockRepository     repository.IStockRepository
	TransactionProducer repository.ITransactionProducer

	TransactionFeedService service.ITransactionFeedService
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

	return &Dependency{
		StockRepository:        stockRepo,
		TransactionProducer:    trxProducer,
		TransactionFeedService: trxFeedService,
	}
}
