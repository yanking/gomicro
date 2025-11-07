// Package database provides support for multiple database instances including MySQL, Redis and MongoDB.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// mongoInstances stores multiple mongo instances
	mongoInstances = make(map[string]*mongo.Client)
	mongoMu        sync.RWMutex
)

// MongoDBOptions defines options for MongoDB connection.
type MongoDBOptions struct {
	// Instance is the name of the MongoDB instance
	Instance string
	// URI is the MongoDB connection URI
	URI string
	// ConnectTimeout is the timeout for establishing connection
	ConnectTimeout time.Duration
	// MaxPoolSize is the maximum number of connections in the connection pool
	MaxPoolSize uint64
	// Logger is the slog logger for MongoDB operations
	Logger *slog.Logger
}

// InitMongoDB initializes a single MongoDB instance.
func InitMongoDB(opts *MongoDBOptions) (*mongo.Client, error) {
	return InitMongoDBWithContext(context.Background(), opts)
}

// InitMongoDBWithContext initializes MongoDB with context.
func InitMongoDBWithContext(ctx context.Context, opts *MongoDBOptions) (*mongo.Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("MongoDB options is nil")
	}

	// Set default timeout if not specified
	connectTimeout := opts.ConnectTimeout
	if connectTimeout == 0 {
		connectTimeout = 10 * time.Second
	}

	// Create client options
	clientOptions := options.Client().
		ApplyURI(opts.URI).
		SetConnectTimeout(connectTimeout).
		SetMaxPoolSize(opts.MaxPoolSize)

	// Set up logger if provided
	if opts.Logger != nil {
		loggerOptions := options.Logger()
		loggerOptions.SetSink(newMongoLogger(opts.Logger))
		clientOptions.SetLoggerOptions(loggerOptions)
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Store the instance
	mongoMu.Lock()
	mongoInstances[opts.Instance] = client
	mongoMu.Unlock()

	return client, nil
}

// InitMongoDBs initializes multiple MongoDB instances.
func InitMongoDBs(opts []*MongoDBOptions) error {
	for _, opt := range opts {
		if _, err := InitMongoDBWithContext(context.Background(), opt); err != nil {
			return fmt.Errorf("failed to initialize MongoDB instance '%s': %w", opt.Instance, err)
		}
	}
	return nil
}

// GetMongoDB returns a MongoDB instance by name.
// If no name is provided or name is empty, it returns the default instance (first one).
func GetMongoDB(instances ...string) *mongo.Client {
	mongoMu.RLock()
	defer mongoMu.RUnlock()

	instance := "default"
	if len(instances) > 0 && instances[0] != "" {
		instance = instances[0]
	}

	if client, exists := mongoInstances[instance]; exists {
		return client
	}

	// Return the first available instance as default
	for _, client := range mongoInstances {
		return client
	}

	return nil
}

// GetMongoDBInstances returns all MongoDB instance names.
func GetMongoDBInstances() []string {
	mongoMu.RLock()
	defer mongoMu.RUnlock()

	instances := make([]string, 0, len(mongoInstances))
	for name := range mongoInstances {
		instances = append(instances, name)
	}
	return instances
}

// CloseMongoDB closes specified MongoDB instances.
// If no instances are specified, all instances will be closed.
func CloseMongoDB(ctx context.Context, instances ...string) error {
	mongoMu.Lock()
	defer mongoMu.Unlock()

	// If no instances specified, close all
	if len(instances) == 0 {
		for name, client := range mongoInstances {
			if err := client.Disconnect(ctx); err != nil {
				return fmt.Errorf("failed to close MongoDB instance '%s': %w", name, err)
			}
			delete(mongoInstances, name)
		}
		return nil
	}

	// Close specified instances
	for _, instance := range instances {
		if client, exists := mongoInstances[instance]; exists {
			if err := client.Disconnect(ctx); err != nil {
				return fmt.Errorf("failed to close MongoDB instance '%s': %w", instance, err)
			}
			delete(mongoInstances, instance)
		}
	}
	return nil
}
