package application

import (
	"context"
	"stockbit-challenge/adapter/kafka/consumer"
	"stockbit-challenge/adapter/redis"
	"stockbit-challenge/config"
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
	}

	app.Config = config.LoadConfiguration()
	app.Consumer = consumer.NewConsumer(ctx, app.Config.Kafka.Consumer)
	app.Redis = redis.NewRedis(app.Config.Redis)

	return app
}

func (app *App) Close() {}
