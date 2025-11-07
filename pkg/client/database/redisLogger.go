package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ redis.Hook = (*redisLogger)(nil)

type redisLogger struct {
	logger *slog.Logger
}

func (r *redisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// Use provided logger or default logger
		start := time.Now()
		conn, err := next(ctx, network, addr)
		duration := time.Since(start)

		if err != nil {
			r.logger.Error("Redis dial failed",
				slog.String("network", network),
				slog.String("addr", addr),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()))
		} else {
			r.logger.Info("Redis dial success",
				slog.String("network", network),
				slog.String("addr", addr),
				slog.Duration("duration", duration))
		}

		return conn, err
	}
}

func (r *redisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {

		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start)

		if err != nil && !errors.Is(err, redis.Nil) {
			r.logger.Error("Redis command failed",
				slog.String("command", cmd.Name()),
				slog.String("args", fmt.Sprintf("%v", cmd.Args())),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()))
		} else {
			r.logger.Info("Redis command executed",
				slog.String("command", cmd.Name()),
				slog.String("args", fmt.Sprintf("%v", cmd.Args())),
				slog.Duration("duration", duration))
		}

		return err
	}
}

func (r *redisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {

		start := time.Now()
		err := next(ctx, cmds)
		duration := time.Since(start)

		cmdNames := make([]string, len(cmds))
		for i, cmd := range cmds {
			cmdNames[i] = cmd.Name()
		}

		if err != nil && !errors.Is(err, redis.Nil) {
			r.logger.Error("Redis pipeline failed",
				slog.String("commands", fmt.Sprintf("%v", cmdNames)),
				slog.Int("count", len(cmds)),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()))
		} else {
			r.logger.Info("Redis pipeline executed",
				slog.String("commands", fmt.Sprintf("%v", cmdNames)),
				slog.Int("count", len(cmds)),
				slog.Duration("duration", duration))
		}

		return err
	}
}
