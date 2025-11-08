// Package mq provides asynq client and server implementations.
package mq

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// AsynqOptions defines options for asynq.
type AsynqOptions struct {
	// Instance is the name of the asynq instance
	Instance string `mapstructure:"instance"`

	// Redis client options
	Redis redis.UniversalClient `mapstructure:"-"`

	// Logger is the logger to use for asynq
	Logger *slog.Logger `mapstructure:"-"`

	// Concurrency is the maximum number of concurrent workers
	Concurrency int `mapstructure:"concurrency"`

	// Queues is a list of queues to process
	Queues map[string]int `mapstructure:"queues"`

	// RedisAddr is the address of the Redis server (for single node)
	RedisAddr string `mapstructure:"redis_addr"`

	// RedisAddrs is the addresses of the Redis servers (for cluster or sentinel)
	RedisAddrs []string `mapstructure:"redis_addrs"`

	// RedisDB is the Redis database number (for single node)
	RedisDB int `mapstructure:"redis_db"`

	// RedisUsername is the Redis username
	RedisUsername string `mapstructure:"redis_username"`

	// RedisPassword is the Redis password
	RedisPassword string `mapstructure:"redis_password"`

	// MasterName is the Redis sentinel master name
	MasterName string `mapstructure:"master_name"`

	// SentinelUsername is the Redis sentinel username
	SentinelUsername string `mapstructure:"sentinel_username"`

	// SentinelPassword is the Redis sentinel password
	SentinelPassword string `mapstructure:"sentinel_password"`

	// SchedulerOpts contains scheduler options
	SchedulerOpts *asynq.SchedulerOpts `mapstructure:"-"`
}

// asynqInstances stores multiple asynq instances
var (
	asynqInstances       = make(map[string]*asynq.Client)
	asynqServerInstances = make(map[string]*asynq.Server)
	schedulerInstances   = make(map[string]*asynq.Scheduler)
	asynqMu              sync.RWMutex
)

// InitAsynq initializes a single asynq instance.
func InitAsynq(opts *AsynqOptions) (*asynq.Client, error) {
	if opts == nil {
		return nil, errors.New("asynq options is nil")
	}

	// Create asynq client
	redisOpt, err := createRedisConnOpt(opts)
	if err != nil {
		return nil, err
	}

	client := asynq.NewClient(redisOpt)

	asynqMu.Lock()
	asynqInstances[opts.Instance] = client
	asynqMu.Unlock()

	return client, nil
}

// GetAsynq returns an asynq instance by name.
func GetAsynq(name string) *asynq.Client {
	asynqMu.RLock()
	defer asynqMu.RUnlock()
	return asynqInstances[name]
}

// InitAsynqServer initializes an asynq server instance.
func InitAsynqServer(opts *AsynqOptions) (*asynq.Server, error) {
	if opts == nil {
		return nil, errors.New("asynq options is nil")
	}

	// Create asynq server config
	serverOpts := asynq.Config{
		Concurrency: opts.Concurrency,
		Queues:      opts.Queues,
	}

	// Create asynq server
	redisOpt, err := createRedisConnOpt(opts)
	if err != nil {
		return nil, err
	}

	server := asynq.NewServer(redisOpt, serverOpts)

	asynqMu.Lock()
	asynqServerInstances[opts.Instance] = server
	asynqMu.Unlock()

	return server, nil
}

// GetAsynqServer returns an asynq server instance by name.
func GetAsynqServer(name string) *asynq.Server {
	asynqMu.RLock()
	defer asynqMu.RUnlock()
	return asynqServerInstances[name]
}

// InitScheduler initializes an asynq scheduler instance.
func InitScheduler(opts *AsynqOptions) (*asynq.Scheduler, error) {
	if opts == nil {
		return nil, errors.New("asynq options is nil")
	}

	// Create asynq scheduler
	redisOpt, err := createRedisConnOpt(opts)
	if err != nil {
		return nil, err
	}

	scheduler := asynq.NewScheduler(redisOpt, opts.SchedulerOpts)

	asynqMu.Lock()
	schedulerInstances[opts.Instance] = scheduler
	asynqMu.Unlock()

	return scheduler, nil
}

// GetScheduler returns an asynq scheduler instance by name.
func GetScheduler(name string) *asynq.Scheduler {
	asynqMu.RLock()
	defer asynqMu.RUnlock()
	return schedulerInstances[name]
}
