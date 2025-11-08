// Package mq provides support for message queue operations including Kafka.
package mq

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

var (
	// kafkaProducers stores multiple kafka producer instances
	kafkaProducers = make(map[string]sarama.SyncProducer)
	// kafkaConsumers stores multiple kafka consumer instances
	kafkaConsumers = make(map[string]sarama.Consumer)
	// kafkaProducerMu protects kafkaProducers
	kafkaProducerMu sync.RWMutex
	// kafkaConsumerMu protects kafkaConsumers
	kafkaConsumerMu sync.RWMutex
)

// KafkaOptions defines options for Kafka connection.
type KafkaOptions struct {
	// Instance is the name of the Kafka instance
	Instance string
	// Brokers is a list of Kafka broker addresses
	Brokers []string
	// Version is the Kafka version (default: "2.1.0")
	Version string
	// Producer configuration
	Producer *KafkaProducerOptions
	// Consumer configuration
	Consumer *KafkaConsumerOptions
	// Logger is the slog logger for Kafka operations
	Logger *slog.Logger
}

// KafkaProducerOptions defines options for Kafka producer.
type KafkaProducerOptions struct {
	// MaxRetries is the maximum number of retries for sending a message
	MaxRetries int
	// RetryBackoff is the time to wait between retries
	RetryBackoff time.Duration
	// RequiredAcks is the number of acks required (default: WaitForLocal)
	RequiredAcks sarama.RequiredAcks
	// Timeout is the maximum time to wait for a response
	Timeout time.Duration
}

// KafkaConsumerOptions defines options for Kafka consumer.
type KafkaConsumerOptions struct {
	// GroupID is the consumer group ID
	GroupID string
	// OffsetInitial is the initial offset position (default: Newest)
	OffsetInitial int64
	// Timeout is the maximum time to wait for a message
	Timeout time.Duration
}

// InitKafka initializes a single Kafka instance for both producer and consumer.
func InitKafka(opts *KafkaOptions) error {
	if opts == nil {
		return fmt.Errorf("Kafka options is nil")
	}

	kafkaVersion := sarama.V2_1_0_0
	if opts.Version != "" {
		version, err := sarama.ParseKafkaVersion(opts.Version)
		if err != nil {
			return fmt.Errorf("failed to parse Kafka version: %w", err)
		}
		kafkaVersion = version
	}

	// Initialize producer if configured
	if opts.Producer != nil {
		producer, err := initKafkaProducer(opts.Brokers, kafkaVersion, opts.Producer, opts.Logger)
		if err != nil {
			return fmt.Errorf("failed to initialize Kafka producer for instance '%s': %w", opts.Instance, err)
		}

		kafkaProducerMu.Lock()
		kafkaProducers[opts.Instance] = producer
		kafkaProducerMu.Unlock()
	}

	// Initialize consumer if configured
	if opts.Consumer != nil {
		consumer, err := initKafkaConsumer(opts.Brokers, kafkaVersion, opts.Consumer, opts.Logger)
		if err != nil {
			return fmt.Errorf("failed to initialize Kafka consumer for instance '%s': %w", opts.Instance, err)
		}

		kafkaConsumerMu.Lock()
		kafkaConsumers[opts.Instance] = consumer
		kafkaConsumerMu.Unlock()
	}

	return nil
}

// initKafkaProducer creates a Kafka producer instance.
func initKafkaProducer(brokers []string, version sarama.KafkaVersion,
	opts *KafkaProducerOptions, logger *slog.Logger) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = version

	// Producer configuration
	if opts.MaxRetries > 0 {
		config.Producer.Retry.Max = opts.MaxRetries
	}
	if opts.RetryBackoff > 0 {
		config.Producer.Retry.Backoff = opts.RetryBackoff
	}
	if opts.RequiredAcks != 0 {
		config.Producer.RequiredAcks = opts.RequiredAcks
	} else {
		config.Producer.RequiredAcks = sarama.WaitForLocal
	}
	if opts.Timeout > 0 {
		config.Producer.Timeout = opts.Timeout
	} else {
		config.Producer.Timeout = 10 * time.Second
	}

	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Set up logger if provided
	_ = logger // Placeholder for future logger implementation

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return producer, nil
}

// initKafkaConsumer creates a Kafka consumer instance.
func initKafkaConsumer(brokers []string, version sarama.KafkaVersion,
	opts *KafkaConsumerOptions, logger *slog.Logger) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Version = version

	// Consumer configuration
	if opts.OffsetInitial != 0 {
		config.Consumer.Offsets.Initial = opts.OffsetInitial
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	if opts.Timeout > 0 {
		config.Consumer.MaxWaitTime = opts.Timeout
	} else {
		config.Consumer.MaxWaitTime = 250 * time.Millisecond
	}

	// Set up logger if provided
	_ = logger // Placeholder for future logger implementation

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return consumer, nil
}

// InitKafkas initializes multiple Kafka instances.
func InitKafkas(opts []*KafkaOptions) error {
	for _, opt := range opts {
		if err := InitKafka(opt); err != nil {
			return fmt.Errorf("failed to initialize Kafka instance '%s': %w", opt.Instance, err)
		}
	}
	return nil
}

// GetKafkaProducer returns a Kafka producer instance by name.
// If no name is provided or name is empty, it returns the default instance (first one).
func GetKafkaProducer(instances ...string) sarama.SyncProducer {
	kafkaProducerMu.RLock()
	defer kafkaProducerMu.RUnlock()

	instance := "default"
	if len(instances) > 0 && instances[0] != "" {
		instance = instances[0]
	}

	if producer, exists := kafkaProducers[instance]; exists {
		return producer
	}

	// Return the first available instance as default
	for _, producer := range kafkaProducers {
		return producer
	}

	return nil
}

// GetKafkaConsumer returns a Kafka consumer instance by name.
// If no name is provided or name is empty, it returns the default instance (first one).
func GetKafkaConsumer(instances ...string) sarama.Consumer {
	kafkaConsumerMu.RLock()
	defer kafkaConsumerMu.RUnlock()

	instance := "default"
	if len(instances) > 0 && instances[0] != "" {
		instance = instances[0]
	}

	if consumer, exists := kafkaConsumers[instance]; exists {
		return consumer
	}

	// Return the first available instance as default
	for _, consumer := range kafkaConsumers {
		return consumer
	}

	return nil
}

// GetKafkaProducerInstances returns all Kafka producer instance names.
func GetKafkaProducerInstances() []string {
	kafkaProducerMu.RLock()
	defer kafkaProducerMu.RUnlock()

	instances := make([]string, 0, len(kafkaProducers))
	for name := range kafkaProducers {
		instances = append(instances, name)
	}
	return instances
}

// GetKafkaConsumerInstances returns all Kafka consumer instance names.
func GetKafkaConsumerInstances() []string {
	kafkaConsumerMu.RLock()
	defer kafkaConsumerMu.RUnlock()

	instances := make([]string, 0, len(kafkaConsumers))
	for name := range kafkaConsumers {
		instances = append(instances, name)
	}
	return instances
}

// closeKafkaComponents closes specified Kafka components.
// If no instances are specified, all instances will be closed.
func closeKafkaComponents[T any](mu *sync.RWMutex,
	components map[string]T,
	closeFunc func(T) error,
	componentType string,
	instances ...string) error {
	mu.Lock()
	defer mu.Unlock()

	// If no instances specified, close all
	if len(instances) == 0 {
		for name, component := range components {
			if err := closeFunc(component); err != nil {
				return fmt.Errorf("failed to close Kafka %s instance '%s': %w", componentType, name, err)
			}
			delete(components, name)
		}
		return nil
	}

	// Close specified instances
	for _, instance := range instances {
		if component, exists := components[instance]; exists {
			if err := closeFunc(component); err != nil {
				return fmt.Errorf("failed to close Kafka %s instance '%s': %w", componentType, instance, err)
			}
			delete(components, instance)
		}
	}
	return nil
}

// closeKafkaProducers closes specified Kafka producer instances.
// If no instances are specified, all producer instances will be closed.
func closeKafkaProducers(_ context.Context, instances ...string) error {
	return closeKafkaComponents(&kafkaProducerMu, kafkaProducers,
		func(p sarama.SyncProducer) error { return p.Close() },
		"producer", instances...)
}

// closeKafkaConsumers closes specified Kafka consumer instances.
// If no instances are specified, all consumer instances will be closed.
func closeKafkaConsumers(_ context.Context, instances ...string) error {
	return closeKafkaComponents(&kafkaConsumerMu, kafkaConsumers,
		func(c sarama.Consumer) error { return c.Close() },
		"consumer", instances...)
}

// CloseKafkaProducer closes specified Kafka producer instances.
// If no instances are specified, all producer instances will be closed.
func CloseKafkaProducer(ctx context.Context, instances ...string) error {
	return closeKafkaProducers(ctx, instances...)
}

// CloseKafkaConsumer closes specified Kafka consumer instances.
// If no instances are specified, all consumer instances will be closed.
func CloseKafkaConsumer(ctx context.Context, instances ...string) error {
	return closeKafkaConsumers(ctx, instances...)
}

// SendMessage sends a message to a Kafka topic using the specified producer instance.
func SendMessage(instance, topic string, key, value []byte) (partition int32, offset int64, err error) {
	producer := GetKafkaProducer(instance)
	if producer == nil {
		return 0, 0, fmt.Errorf("Kafka producer instance '%s' not found", instance)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err = producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send message to topic '%s': %w", topic, err)
	}

	return partition, offset, nil
}

// ConsumeMessages consumes messages from a Kafka topic using the specified consumer instance.
// This function returns a channel that receives messages and should be run in a goroutine.
func ConsumeMessages(ctx context.Context, instance, topic string, partition int32) (<-chan *sarama.ConsumerMessage, <-chan error) {
	messageChan := make(chan *sarama.ConsumerMessage)
	errorChan := make(chan error)

	go func() {
		defer close(messageChan)
		defer close(errorChan)

		consumer := GetKafkaConsumer(instance)
		if consumer == nil {
			errorChan <- fmt.Errorf("Kafka consumer instance '%s' not found", instance)
			return
		}

		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			errorChan <- fmt.Errorf("failed to consume partition %d of topic '%s': %w", partition, topic, err)
			return
		}
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				// Log the error or handle it as appropriate for your application
				_ = err // For now, just acknowledge the error
			}
		}()

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				select {
				case messageChan <- msg:
				case <-ctx.Done():
					return
				}
			case err := <-partitionConsumer.Errors():
				select {
				case errorChan <- err:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return messageChan, errorChan
}
