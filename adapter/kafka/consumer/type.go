package consumer

import (
	"context"
	"stockbit-challenge/adapter/kafka"
)

//go:generate mockgen -destination=mock/consumer.go -package=mock . IConsumer

type Consumer struct {
	Config *Configuration
}

type IConsumer interface {
	Close() error

	Consume(ctx context.Context, topic string) (*kafka.Message, error)

	Fetch(ctx context.Context, topic string) (*kafka.Message, error)
	Commit(ctx context.Context, topic string, message *kafka.Message) error
}
