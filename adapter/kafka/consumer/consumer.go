package consumer

import (
	"context"
	"errors"
	"sync"

	kafkaGo "github.com/segmentio/kafka-go"

	"stockbit-challenge/adapter/kafka"
)

const (
	LastOffset  int64 = -1 // The most recent offset available for a partition.
	FirstOffset int64 = -2 // The least recent offset available for a partition.
)

var (
	consumers sync.Map
)

// NewConsumer will create Consumer instance based on context and Kafka configuration
func NewConsumer(config *Configuration, opts ...Option) IConsumer {
	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	consumer := &Consumer{
		Config: config,
	}

	if options.groupID != nil {
		consumer.Config.GroupID = *options.groupID
	}
	if options.startOffset != nil {
		consumer.Config.StartOffset = *options.startOffset
	}

	return consumer
}

// Consume will get next message and commit it immediately
func (k *Consumer) Consume(ctx context.Context, topic string) (*kafka.Message, error) {
	if topic == "" {
		return nil, errors.New("empty topic")
	}

	consumer := k.getConsumer(topic)

	if msg, err := consumer.ReadMessage(ctx); err != nil {
		return nil, err
	} else {
		headers := make(kafka.Header)

		for _, header := range msg.Headers {
			headers[header.Key] = header.Value
		}

		return &kafka.Message{
			Partition: msg.Partition,
			Offset:    msg.Offset,

			Key:   msg.Key,
			Value: msg.Value,

			Headers: headers,
		}, nil
	}
}

// Fetch will get next message but will not commit it
// see also: Commit(topic string, message Message)
func (k *Consumer) Fetch(ctx context.Context, topic string) (*kafka.Message, error) {
	if topic == "" {
		return nil, errors.New("empty topic")
	}

	consumer := k.getConsumer(topic)

	if msg, err := consumer.FetchMessage(ctx); err != nil {
		return nil, err
	} else {
		headers := make(kafka.Header)

		for _, header := range msg.Headers {
			headers[header.Key] = header.Value
		}

		return &kafka.Message{
			Partition: msg.Partition,
			Offset:    msg.Offset,

			Key:   msg.Key,
			Value: msg.Value,

			Headers: headers,
		}, nil
	}
}

// Commit will commit message. only need topic, partition and offset.
func (k *Consumer) Commit(ctx context.Context, topic string, message *kafka.Message) error {
	if topic == "" {
		return errors.New("empty topic")
	}

	consumer := k.getConsumer(topic)

	msg := kafkaGo.Message{ // only needs these three
		Topic: topic,

		Partition: message.Partition,
		Offset:    message.Offset,
	}

	return consumer.CommitMessages(ctx, msg)
}

// getConsumer get list of consumer for a specific topic
func (k *Consumer) getConsumer(topic string /*, partition int*/) *kafkaGo.Reader {
	var consumer *kafkaGo.Reader

	if v, ok := consumers.Load(topic); ok {
		consumer = v.(*kafkaGo.Reader)
	} else {
		consumer = kafkaGo.NewReader(
			kafkaGo.ReaderConfig{
				Brokers: k.Config.Brokers,

				GroupID:     k.Config.GroupID,
				GroupTopics: k.Config.GroupTopics,

				// The topic to read messages from.
				Topic: topic,

				// Partition to read messages from.  Either Partition or GroupID may
				// be assigned, but not both
				// Partition: partition,

				QueueCapacity: k.Config.QueueCapacity,

				MinBytes: k.Config.MinBytes,
				MaxBytes: k.Config.MaxBytes,

				MaxWait: k.Config.MaxWait,

				ReadLagInterval:        k.Config.ReadLagInterval,
				HeartbeatInterval:      k.Config.HeartbeatInterval,
				CommitInterval:         k.Config.CommitInterval,
				PartitionWatchInterval: k.Config.PartitionWatchInterval,

				WatchPartitionChanges: k.Config.WatchPartitionChanges,

				SessionTimeout:   k.Config.SessionTimeout,
				RebalanceTimeout: k.Config.RebalanceTimeout,

				JoinGroupBackoff: k.Config.JoinGroupBackoff,
				RetentionTime:    k.Config.RetentionTime,

				StartOffset: k.Config.StartOffset,

				ReadBackoffMin: k.Config.ReadBackoffMax,
				ReadBackoffMax: k.Config.ReadBackoffMin,

				MaxAttempts: k.Config.MaxAttempts,
			},
		)

		// Add consumer to list of known consumers by topic
		consumers.Store(topic, consumer)
	}

	return consumer
}

// Close all known consumers
func (*Consumer) Close() error {
	var err error

	var wg sync.WaitGroup
	defer wg.Wait()

	consumers.Range(
		func(topic, consumer interface{}) bool {
			if consumer != nil {
				wg.Add(1)

				go func() {
					// Close the writer
					err = consumer.(*kafkaGo.Reader).Close()

					wg.Done()
				}()
			}

			// Remove consumer from list of consumers
			consumers.Delete(topic)

			return true
		})

	return err
}
