// Package database provides support for multiple database instances including MySQL and Redis.
package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

// redisInstances stores multiple redis instances
var (
	redisInstances = make(map[string]redis.UniversalClient)
	redisMu        sync.RWMutex
)

// RedisOptions defines configuration for Redis (single instance or cluster).
// If Addrs length > 1, cluster mode will be used automatically.
// Otherwise, the first address (Addrs[0]) will be used as single-node.
type RedisOptions struct {
	// Instance is the name of the redis instance
	Instance string
	// Addrs is a list of redis addresses. Provide at least one address.
	// If multiple addresses are provided, cluster mode will be used automatically.
	Addrs         []string
	Username      string
	Password      string
	DB            int
	PoolSize      int
	MinIdleConns  int
	DialTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	EnableTLS     bool
	TLSSkipVerify bool
	TLSCAFile     string
	TLSCertFile   string
	TLSKeyFile    string
	TLSServerName string
	// Logger is the slog logger for Redis operations
	Logger *slog.Logger
}

// InitRedis initializes a single redis instance.
func InitRedis(opts *RedisOptions) (redis.UniversalClient, error) {
	return InitRedisWithContext(context.Background(), opts)
}

// InitRedisWithContext initializes with context.
// It supports both single Redis instance and Redis cluster modes.
// When Addrs contains multiple addresses, cluster mode is used automatically.
func InitRedisWithContext(ctx context.Context, opts *RedisOptions) (redis.UniversalClient, error) {
	if opts == nil {
		return nil, errors.New("redis options is nil")
	}
	addrs := opts.Addrs
	if len(addrs) == 0 {
		return nil, errors.New("redis addrs is empty")
	}

	client, err := createRedisClient(opts, addrs)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	// Add logging hook with provided logger or default logger
	if opts.Logger != nil {
		logger := opts.Logger.With(slog.String("component", "redis"))
		client.AddHook(&redisLogger{logger: logger})
	}

	// Store the instance
	redisMu.Lock()
	redisInstances[opts.Instance] = client
	redisMu.Unlock()

	return client, nil
}

// createRedisClient creates a redis client with the given options
func createRedisClient(opts *RedisOptions, addrs []string) (redis.UniversalClient, error) {
	// TLS
	var tlsConfig *tls.Config
	var err error
	if opts.EnableTLS {
		tlsConfig, err = createTLSConfig(opts, addrs)
		if err != nil {
			return nil, err
		}
	}

	uo := &redis.UniversalOptions{
		Addrs:        addrs,
		Username:     opts.Username,
		Password:     opts.Password,
		DB:           opts.DB,
		PoolSize:     opts.PoolSize,
		MinIdleConns: opts.MinIdleConns,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		TLSConfig:    tlsConfig,
		// Explicitly disable maintenance notifications
		// This prevents the client from sending CLIENT MAINT_NOTIFICATIONS ON
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	}

	return redis.NewUniversalClient(uo), nil
}

// createTLSConfig creates a TLS configuration for Redis client
func createTLSConfig(opts *RedisOptions, addrs []string) (*tls.Config, error) {
	cfg := &tls.Config{MinVersion: tls.VersionTLS12}
	if opts.TLSSkipVerify {
		cfg.InsecureSkipVerify = true
	}
	if opts.TLSServerName != "" {
		cfg.ServerName = opts.TLSServerName
	}
	if opts.TLSCAFile != "" {
		caBytes, err := os.ReadFile(opts.TLSCAFile)
		if err != nil {
			return nil, fmt.Errorf("read CA file: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caBytes) {
			return nil, errors.New("append CA cert failed")
		}
		cfg.RootCAs = pool
	}
	if opts.TLSCertFile != "" && opts.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(opts.TLSCertFile, opts.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("load client cert/key: %w", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}
	if !cfg.InsecureSkipVerify && cfg.ServerName == "" { // derive from first addr
		if host, _, err := net.SplitHostPort(addrs[0]); err == nil && host != "" {
			cfg.ServerName = host
		}
	}
	return cfg, nil
}

// InitRedises initializes multiple redis instances.
func InitRedises(opts []*RedisOptions) error {
	for _, opt := range opts {
		if _, err := InitRedisWithContext(context.Background(), opt); err != nil {
			return fmt.Errorf("failed to initialize Redis instance '%s': %w", opt.Instance, err)
		}
	}
	return nil
}

// GetRedis returns a redis instance by name.
// If no name is provided or name is empty, it returns the default instance (first one).
func GetRedis(instances ...string) redis.UniversalClient {
	redisMu.RLock()
	defer redisMu.RUnlock()

	instance := "default"
	if len(instances) > 0 && instances[0] != "" {
		instance = instances[0]
	}

	if client, exists := redisInstances[instance]; exists {
		return client
	}

	// Return the first available instance as default
	for _, client := range redisInstances {
		return client
	}

	return nil
}

// GetRedisInstances returns all redis instance names.
func GetRedisInstances() []string {
	redisMu.RLock()
	defer redisMu.RUnlock()

	instances := make([]string, 0, len(redisInstances))
	for name := range redisInstances {
		instances = append(instances, name)
	}
	return instances
}

// CloseRedis closes specified redis instances.
// If no instances are specified, all instances will be closed.
func CloseRedis(_ context.Context, instances ...string) error {
	redisMu.Lock()
	defer redisMu.Unlock()

	// If no instances specified, close all
	if len(instances) == 0 {
		for name, client := range redisInstances {
			if err := client.Close(); err != nil {
				return fmt.Errorf("failed to close Redis instance '%s': %w", name, err)
			}
			delete(redisInstances, name)
		}
		return nil
	}

	// Close specified instances
	for _, instance := range instances {
		if client, exists := redisInstances[instance]; exists {
			if err := client.Close(); err != nil {
				return fmt.Errorf("failed to close Redis instance '%s': %w", instance, err)
			}
			delete(redisInstances, instance)
		}
	}
	return nil
}
