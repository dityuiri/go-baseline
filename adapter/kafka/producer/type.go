package producer

//go:generate mockgen -destination=mock/producer.go -package=mock . IProducer

import (
	"context"

	"stockbit-challenge/adapter/kafka"
)

type Producer struct {
	Context context.Context
	Config  *Configuration
}

type IProducer interface {
	Close() error

	Produce(topic string, messages ...*kafka.Message) error
}
