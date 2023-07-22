package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"stockbit-challenge/adapter/kafka"
	"stockbit-challenge/adapter/kafka/producer"
	"stockbit-challenge/config"
	"stockbit-challenge/model"
)

//go:generate mockgen -package=repository_mock -destination=../mock/repository/transaction_producer.go . ITransactionProducer

type (
	ITransactionProducer interface {
		ProduceTrx(ctx context.Context, trx model.Transaction) error
		ProduceTrxDLQ(ctx context.Context, trx model.Transaction, err error) error
	}

	TransactionProducer struct {
		Producer    producer.IProducer
		KafkaConfig *config.Kafka
	}
)

func (p *TransactionProducer) ProduceTrx(ctx context.Context, trx model.Transaction) error {
	msg := p.constructMessage(trx)
	return p.Producer.Produce(ctx, p.KafkaConfig.ProducerTopics["transaction"], msg)
}

func (p *TransactionProducer) ProduceTrxDLQ(ctx context.Context, trx model.Transaction, err error) error {
	trx.Error = err.Error()
	msg := p.constructMessage(trx)
	return p.Producer.Produce(ctx, p.KafkaConfig.ProducerTopics["transaction_dlq"], msg)
}

func (*TransactionProducer) constructMessage(trx model.Transaction) *kafka.Message {
	var message *kafka.Message
	data, _ := json.Marshal(trx)
	message = &kafka.Message{
		Value: data,
		Headers: map[string][]byte{
			"message_id": []byte(fmt.Sprintf("%s-%v", trx.StockCode, trx.OrderBook)),
		},
	}

	return message
}
