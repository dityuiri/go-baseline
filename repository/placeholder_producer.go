package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dityuiri/go-baseline/adapter/kafka"
	"github.com/dityuiri/go-baseline/adapter/kafka/producer"
	"github.com/dityuiri/go-baseline/config"
	"github.com/dityuiri/go-baseline/model"
)

//go:generate mockgen -package=repository_mock -destination=../mock/repository/placeholder_producer.go . IPlaceholderProducer

type (
	IPlaceholderProducer interface {
		ProducePlaceholderRecord(ctx context.Context, placeholderMsg model.PlaceholderMessage) error
	}

	PlaceholderProducer struct {
		Producer    producer.IProducer
		KafkaConfig *config.Kafka
	}
)

func (p *PlaceholderProducer) ProducePlaceholderRecord(ctx context.Context, placeholderMsg model.PlaceholderMessage) error {
	msg := p.constructMessage(placeholderMsg)
	return p.Producer.Produce(ctx, p.KafkaConfig.ProducerTopics["placeholder"], msg)
}

func (*PlaceholderProducer) constructMessage(placeholderMsg model.PlaceholderMessage) *kafka.Message {
	var message *kafka.Message
	data, _ := json.Marshal(placeholderMsg)
	message = &kafka.Message{
		Value: data,
		Headers: map[string][]byte{
			"message_id": []byte(fmt.Sprintf("%s", placeholderMsg.ID)),
		},
	}

	return message
}
