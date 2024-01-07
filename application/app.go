package application

import (
	"context"
	"github.com/dityuiri/go-baseline/adapter/db"

	"github.com/dityuiri/go-baseline/adapter/kafka/consumer"
	"github.com/dityuiri/go-baseline/adapter/logger"
	"github.com/dityuiri/go-baseline/adapter/redis"
	"github.com/dityuiri/go-baseline/config"
)

type App struct {
	Context  context.Context
	Config   *config.Configuration
	Consumer consumer.IConsumer
	Redis    redis.IRedis
	Logger   logger.ILogger
	DB       db.IDatabase
}

func SetupApplication(ctx context.Context) (*App, error) {
	app := &App{
		Context: ctx,
		Config:  config.LoadConfiguration(),
	}

	app.Consumer = consumer.NewConsumer(app.Config.Kafka.Consumer)
	app.Redis = redis.NewRedis(app.Config.Redis)

	loggerInstance, err := logger.NewLogger(logger.WithAppName(app.Config.AppName))
	if err != nil {
		return nil, err
	}

	app.Logger = loggerInstance

	dbInstance, err := db.NewDatabase(ctx, app.Config.Database)
	if err != nil {
		return nil, err
	}

	app.DB = dbInstance

	return app, nil
}

func (app *App) Close() {
	if app.DB != nil {
		_ = app.DB.Close()
	}

	app.Logger.Info("APP SUCCESSFULLY CLOSED")
}
