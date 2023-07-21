package repository

import (
	"encoding/json"

	"stockbit-challenge/adapter/kafka"
	"stockbit-challenge/adapter/kafka/producer"
	"stockbit-challenge/config"
	"stockbit-challenge/model"
)

//go:generate mockgen -package=repository_mock -destination=../mock/repository/transaction_producer.go . ITransactionProducer

type (
	ITransactionProducer interface {
		ProduceTrx(trx model.Transaction) error
		ProduceTrxDLQ(trx model.Transaction, err error) error
	}

	TransactionProducer struct {
		Producer    producer.IProducer
		KafkaConfig *config.Kafka
	}
)

func (p *TransactionProducer) ProduceTrx(trx model.Transaction) error {
	msg := p.constructMessage(trx)
	return p.Producer.Produce(p.KafkaConfig.ProducerTopics["transaction"], msg)
}

func (p *TransactionProducer) ProduceTrxDLQ(trx model.Transaction, err error) error {
	trx.Error = err.Error()
	msg := p.constructMessage(trx)
	return p.Producer.Produce(p.KafkaConfig.ProducerTopics["transaction_dlq"], msg)
}

func (p *TransactionProducer) constructMessage(trx model.Transaction) *kafka.Message {
	var message *kafka.Message
	data, _ := json.Marshal(trx)
	message = &kafka.Message{
		Value: data,
		Headers: map[string][]byte{
			"message_id": []byte(trx.OrderNumber),
		},
	}

	return message
}
