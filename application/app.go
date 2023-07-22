package application

import (
	"stockbit-challenge/adapter/kafka/consumer"
	"stockbit-challenge/adapter/redis"
	"stockbit-challenge/config"
)

type App struct {
	Config   *config.Configuration
	Consumer consumer.IConsumer
	Redis    redis.IRedis
}

func SetupApplication() *App {
	app := &App{
		Config: config.LoadConfiguration(),
	}

	app.Consumer = consumer.NewConsumer(app.Config.Kafka.Consumer)
	app.Redis = redis.NewRedis(app.Config.Redis)

	return app
}
