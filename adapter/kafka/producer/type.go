package producer

//go:generate mockgen -destination=mock/producer.go -package=mock . IProducer

import (
	"context"

	"stockbit-challenge/adapter/kafka"
)

type Producer struct {
	Config *Configuration
}

type IProducer interface {
	Close() error

	Produce(ctx context.Context, topic string, messages ...*kafka.Message) error
}
