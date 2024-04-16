package config

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/dityuiri/go-adapter/client"
	"github.com/dityuiri/go-adapter/db"
	"github.com/dityuiri/go-adapter/kafka/consumer"
	"github.com/dityuiri/go-adapter/kafka/producer"
	"github.com/dityuiri/go-adapter/redis"
)

type (
	Configuration struct {
		AppName    string
		Const      *Constants
		Kafka      *Kafka
		Redis      *redis.Config
		Database   *db.Configuration
		HTTPClient *HttpClient
	}

	Kafka struct {
		Consumer       *consumer.Configuration
		Producer       *producer.Configuration
		ConsumerTopics map[string]string
		ProducerTopics map[string]string
	}

	Constants struct {
		GRPCPort     int
		HTTPPort     int
		ShortTimeout int
	}

	HttpClient struct {
		ClientConfig *client.Configuration
		ProxyURLs    ProxyURLs
	}

	ProxyURLs struct {
		AlphaURL string
	}
)

func LoadConfiguration() *Configuration {
	// Initialize viper
	viper.AutomaticEnv()
	return &Configuration{
		AppName:    viper.GetString("APP_NAME"),
		Const:      loadConstants(),
		Redis:      redis.NewConfig(),
		Kafka:      loadKafkaConfig(),
		Database:   loadDatabaseConfig(),
		HTTPClient: loadHTTPClientConfig(),
	}
}

func loadConstants() *Constants {
	return &Constants{
		GRPCPort:     viper.GetInt("GRPC_PORT"),
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

func loadDatabaseConfig() *db.Configuration {
	return db.NewConfig()
}

func loadHTTPClientConfig() *HttpClient {
	return &HttpClient{
		ClientConfig: client.NewConfig(),
		ProxyURLs: ProxyURLs{
			AlphaURL: viper.GetString("ALPHA_URL"),
		},
	}
}
