package config

import (
	"strings"

	"github.com/spf13/viper"

	"stockbit-challenge/adapter/kafka/consumer"
	"stockbit-challenge/adapter/kafka/producer"
	"stockbit-challenge/adapter/redis"
)

type (
	Configuration struct {
		Const *Constants
		Kafka *Kafka
		Redis *redis.Config
	}

	Kafka struct {
		Consumer       *consumer.Configuration
		Producer       *producer.Configuration
		ConsumerTopics map[string]string
		ProducerTopics map[string]string
	}

	Constants struct {
		HTTPPort     int
		ShortTimeout int
	}
)

func LoadConfiguration() *Configuration {
	return &Configuration{
		Const: loadConstants(),
		Redis: redis.NewConfig(),
		Kafka: loadKafkaConfig(),
	}
}

func loadConstants() *Constants {
	return &Constants{
		HTTPPort:     viper.GetInt("HTTP_PORT"),
		ShortTimeout: viper.GetInt("SHORT_TIMEOUT"),
	}
}

func loadKafkaConfig() *Kafka {
	var (
		consumerTopics       = strings.Split(strings.TrimSpace(viper.GetString("CONSUMER_TOPICS")), ";")
		producerTopics       = strings.Split(strings.TrimSpace(viper.GetString("PRODUCER_TOPICS")), ";")
		mappedConsumerTopics = map[string]string{}
		mappedProducerTopics = map[string]string{}
	)

	for _, topic := range consumerTopics {
		t := strings.Split(strings.TrimSpace(topic), ":")
		if len(t) == 2 {
			mappedConsumerTopics[t[0]] = t[1]
		}
	}

	for _, topic := range producerTopics {
		t := strings.Split(strings.TrimSpace(topic), ":")
		if len(t) == 2 {
			mappedProducerTopics[t[0]] = t[1]
		}
	}

	return &Kafka{
		Consumer: &consumer.Configuration{
			Brokers:     strings.Split(viper.GetString("KAFKA_BROKERS"), ","),
			GroupID:     viper.GetString("KAFKA_GROUP_ID"),
			MinBytes:    10e3,
			MaxBytes:    10e6,
			StartOffset: consumer.LastOffset,
		},
		ConsumerTopics: mappedConsumerTopics,
		Producer: &producer.Configuration{
			Brokers:      strings.Split(viper.GetString("KAFKA_BROKERS"), ","),
			Async:        false,
			RequiredAcks: 1,
		},
		ProducerTopics: mappedProducerTopics,
	}

}
