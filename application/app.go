package application

import (
	"context"

	"github.com/dityuiri/go-baseline/adapter/kafka/consumer"
	"github.com/dityuiri/go-baseline/adapter/redis"
	"github.com/dityuiri/go-baseline/config"
)

type App struct {
	Context  context.Context
	Config   *config.Configuration
	Consumer consumer.IConsumer
	Redis    redis.IRedis
}

func SetupApplication(ctx context.Context) *App {
	app := &App{
		Context: ctx,
		Config:  config.LoadConfiguration(),
	}

	app.Consumer = consumer.NewConsumer(app.Config.Kafka.Consumer)
	app.Redis = redis.NewRedis(app.Config.Redis)

	return app
}
