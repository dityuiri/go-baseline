package producer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	kafkaGo "github.com/segmentio/kafka-go"

	"github.com/dityuiri/go-baseline/adapter/kafka"
)

const (
	MessageIDHeaderName = "X-Message-Id"
)

var (
	producers sync.Map
)

// NewProducer will create Producer instance based on context and Kafka configuration
func NewProducer(config *Configuration) IProducer {
	producer := &Producer{
		Config: config,
	}

	return producer
}

// getProducer get list of producer for a specific topic
func (k *Producer) getProducer(topic string) *kafkaGo.Writer {
	var producer *kafkaGo.Writer

	if v, ok := producers.Load(topic); ok {
		producer = v.(*kafkaGo.Writer)
	} else {
		producer = kafkaGo.NewWriter(
			kafkaGo.WriterConfig{
				Brokers: k.Config.Brokers,

				Topic: topic,

				Balancer: k.Config.Balancer,

				MaxAttempts: k.Config.MaxAttempts,

				BatchSize:    k.Config.BatchSize,
				BatchBytes:   k.Config.BatchBytes,
				BatchTimeout: k.Config.BatchTimeout,

				ReadTimeout:  k.Config.ReadTimeout,
				WriteTimeout: k.Config.WriteTimeout,

				RequiredAcks: k.Config.RequiredAcks,

				Async: k.Config.Async,
			},
		)

		// Add producer to list of known producers by topic
		producers.Store(topic, producer)
	}

	return producer
}

func (*Producer) transformHeaders(header kafka.Header) []kafkaGo.Header {
	var headers []kafkaGo.Header

	hasMessageID := false

	for k, v := range header {
		headers = append(
			headers,

			kafkaGo.Header{
				Key:   k,
				Value: v,
			},
		)

		if k == MessageIDHeaderName {
			hasMessageID = true
		}
	}

	if !hasMessageID {
		headers = append(
			headers,

			kafkaGo.Header{
				Key:   MessageIDHeaderName,
				Value: []byte(uuid.New().String()),
			},
		)
	}

	return headers
}

// Produce do the producing a message to specific topic
func (k *Producer) Produce(ctx context.Context, topic string, messages ...*kafka.Message) error {
	if topic == "" {
		return errors.New("empty topic")
	}

	msgs := []kafkaGo.Message{}

	for _, message := range messages {
		msg := kafkaGo.Message{
			Offset: message.Offset,

			Key: message.Key,

			Headers: k.transformHeaders(message.Headers),
		}

		switch t := message.Value.(type) {
		case []byte:
			msg.Value = message.Value.([]byte)
		case nil:
		default:
			panic(fmt.Sprintf("unknown message value type %v", t))
		}

		msgs = append(msgs, msg)
	}

	producer := k.getProducer(topic)

	// Write message
	return producer.WriteMessages(
		ctx,
		msgs...,
	)
}

// Close all known producers
func (*Producer) Close() error {
	var err error

	var wg sync.WaitGroup
	defer wg.Wait()

	producers.Range(
		func(topic, producer interface{}) bool {
			if producer != nil {
				wg.Add(1)

				go func() {
					// Close the writer
					err = producer.(*kafkaGo.Writer).Close()

					wg.Done()
				}()
			}

			// Remove producer from list of producers
			producers.Delete(topic)

			return true
		},
	)

	return err
}
