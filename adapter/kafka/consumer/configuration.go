package consumer

import (
	"time"

	"github.com/spf13/viper"
)

// Configuration contains the configuration.
type Configuration struct {
	// The list of broker addresses used to connect to the kafka cluster.
	Brokers []string

	// GroupID holds the optional consumer group id.  If GroupID is specified, then
	// Partition should NOT be specified e.g. 0
	GroupID string
	// GroupTopics allows specifying multiple topics, but can only be used in
	// combination with GroupID, as it is a consumer-group feature. As such, if
	// GroupID is set, then either Topic or GroupTopics must be defined.
	GroupTopics []string

	// The capacity of the internal message queue, defaults to 100 if none is
	// set.
	QueueCapacity int

	// MinBytes indicates to the broker the minimum batch size that the consumer
	// will accept. Setting a high minimum when consuming from a low-volume topic
	// may result in delayed delivery when the broker does not have enough data to
	// satisfy the defined minimum.
	//
	// Default: 1
	MinBytes int
	// MaxBytes indicates to the broker the maximum batch size that the consumer
	// will accept. The broker will truncate a message to satisfy this maximum, so
	// choose a value that is high enough for your largest message size.
	//
	// Default: 1MB
	MaxBytes int

	// Maximum amount of time to wait for new data to come when fetching batches
	// of messages from kafka.
	MaxWait time.Duration

	// ReadLagInterval sets the frequency at which the reader lag is updated.
	// Setting this field to a negative value disables lag reporting.
	ReadLagInterval time.Duration

	// HeartbeatInterval sets the optional frequency at which the reader sends the consumer
	// group heartbeat update.
	//
	// Default: 3s
	//
	// Only used when GroupID is set
	HeartbeatInterval time.Duration

	// CommitInterval indicates the interval at which offsets are committed to
	// the broker.  If 0, commits will be handled synchronously.
	//
	// Default: 0
	//
	// Only used when GroupID is set
	CommitInterval time.Duration

	// PartitionWatchInterval indicates how often a reader checks for partition changes.
	// If a reader sees a partition change (such as a partition add) it will rebalance the group
	// picking up new partitions.
	//
	// Default: 5s
	//
	// Only used when GroupID is set and WatchPartitionChanges is set.
	PartitionWatchInterval time.Duration

	// WatchForPartitionChanges is used to inform kafka-go that a consumer group should be
	// polling the brokers and rebalancing if any partition changes happen to the topic.
	WatchPartitionChanges bool

	// SessionTimeout optionally sets the length of time that may pass without a heartbeat
	// before the coordinator considers the consumer dead and initiates a rebalance.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	SessionTimeout time.Duration

	// RebalanceTimeout optionally sets the length of time the coordinator will wait
	// for members to join as part of a rebalance.  For kafka servers under higher
	// load, it may be useful to set this value higher.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	RebalanceTimeout time.Duration

	// JoinGroupBackoff optionally sets the length of time to wait between re-joining
	// the consumer group after an error.
	//
	// Default: 5s
	JoinGroupBackoff time.Duration

	// RetentionTime optionally sets the length of time the consumer group will be saved
	// by the broker
	//
	// Default: 24h
	//
	// Only used when GroupID is set
	RetentionTime time.Duration

	// StartOffset determines from whence the consumer group should begin
	// consuming when it finds a partition without a committed offset.  If
	// non-zero, it must be set to one of FirstOffset or LastOffset.
	//
	// Default: FirstOffset
	//
	// Only used when GroupID is set
	StartOffset int64

	// BackoffDelayMin optionally sets the smallest amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 100ms
	ReadBackoffMin time.Duration

	// BackoffDelayMax optionally sets the maximum amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 1s
	ReadBackoffMax time.Duration

	// Limit of how many attempts will be made before delivering the error.
	//
	// The default is to try 3 times.
	MaxAttempts int
}

// NewConfig returns an instance to the configuration.
func NewConfig() *Configuration {
	return &Configuration{
		Brokers:     viper.GetStringSlice("KAFKA_BROKERS"),
		GroupID:     viper.GetString("KAFKA_GROUP_ID"),
		StartOffset: viper.GetInt64("KAFKA_START_OFFSET"),
	}
}
