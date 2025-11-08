package mq

import (
	"errors"

	"github.com/hibiken/asynq"
)

// createRedisConnOpt creates a Redis connection option based on the provided options.
func createRedisConnOpt(opts *AsynqOptions) (asynq.RedisConnOpt, error) {
	var redisOpt asynq.RedisConnOpt
	if opts.RedisAddr != "" {
		// Single node Redis
		redisOpt = asynq.RedisClientOpt{
			Addr:     opts.RedisAddr,
			DB:       opts.RedisDB,
			Username: opts.RedisUsername,
			Password: opts.RedisPassword,
		}
	} else if len(opts.RedisAddrs) > 0 {
		if opts.MasterName != "" {
			// Redis Sentinel
			redisOpt = asynq.RedisFailoverClientOpt{
				MasterName:       opts.MasterName,
				SentinelAddrs:    opts.RedisAddrs,
				SentinelUsername: opts.SentinelUsername,
				SentinelPassword: opts.SentinelPassword,
				Username:         opts.RedisUsername,
				Password:         opts.RedisPassword,
				DB:               opts.RedisDB,
			}
		} else {
			// Redis Cluster
			redisOpt = asynq.RedisClusterClientOpt{
				Addrs:    opts.RedisAddrs,
				Username: opts.RedisUsername,
				Password: opts.RedisPassword,
			}
		}
	} else if opts.Redis != nil {
		// 如果提供了 Redis 客户端，我们需要从中提取连接信息
		// 由于 redis.UniversalClient 接口没有 Options() 方法，
		// 我们需要在初始化时直接提供连接信息
		return nil, errors.New("redis client cannot be used directly, please provide RedisAddr or RedisAddrs")
	} else {
		return nil, errors.New("either Redis client, RedisAddr or RedisAddrs must be provided")
	}

	return redisOpt, nil
}
