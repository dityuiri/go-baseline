package producer

import (
	"time"

	kafkaGo "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

// Configuration contains the configuration.
type Configuration struct {
	// The list of brokers used to discover the partitions available on the
	// kafka cluster.
	Brokers []string

	// The balancer used to distribute messages across partitions.
	// The default is to use round robin distribution
	// Selection: round-robin(0), least bytes(1), hash(2)
	Balancer kafkaGo.Balancer

	// Limit on how many attempts will be made to deliver a message.
	MaxAttempts int

	// Limit on how many messages will be buffered before being sent to a
	// partition.
	BatchSize int

	// Limit the maximum size of a request in bytes before being sent to
	// a partition.
	BatchBytes int

	// Time limit on how often incomplete message batches will be flushed to
	// kafka.
	BatchTimeout time.Duration

	// Timeout for read operations performed by the Writer.
	ReadTimeout time.Duration

	// Timeout for write operation performed by the Writer.
	WriteTimeout time.Duration

	// Number of acknowledges from partition replicas required before receiving
	// a response to a produce request (default to -1, which means to wait for
	// all replicas).
	RequiredAcks int

	// Setting this flag to true causes the WriteMessages method to never block.
	// It also means that errors are ignored since the caller will not receive
	// the returned value. Use this only if you don't care about guarantees of
	// whether the messages were written to kafka.
	Async bool
}

// NewConfig returns an instance to the configuration.
func NewConfig() *Configuration {
	return &Configuration{
		Brokers: viper.GetStringSlice("KAFKA_BROKERS"),
	}
}
