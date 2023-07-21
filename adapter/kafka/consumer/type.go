package consumer

import (
	"context"

	"stockbit-challenge/adapter/kafka"
)

type Consumer struct {
	Context context.Context
	Config  *Configuration
}

type IConsumer interface {
	Close() error

	Consume(topic string) (*kafka.Message, error)

	Fetch(topic string) (*kafka.Message, error)
	Commit(topic string, message *kafka.Message) error
}
