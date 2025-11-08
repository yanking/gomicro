// Package serverinterceptors provides common gRPC server interceptors.
package serverinterceptors

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor returns a new unary server interceptor that logs incoming requests.
func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		logger.Info("gRPC request started",
			slog.String("method", info.FullMethod),
			slog.String("start_time", startTime.Format(time.RFC3339)),
		)

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC request failed",
				slog.String("method", info.FullMethod),
				slog.String("duration", duration.String()),
				slog.String("error", err.Error()),
				slog.String("code", st.Code().String()),
			)
		} else {
			logger.Info("gRPC request completed",
				slog.String("method", info.FullMethod),
				slog.String("duration", duration.String()),
			)
		}

		return resp, err
	}
}

// LoggingStreamInterceptor returns a new stream server interceptor that logs incoming streams.
func LoggingStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		startTime := time.Now()

		logger.Info("gRPC stream started",
			slog.String("method", info.FullMethod),
			slog.String("start_time", startTime.Format(time.RFC3339)),
		)

		err := handler(srv, ss)

		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC stream failed",
				slog.String("method", info.FullMethod),
				slog.String("duration", duration.String()),
				slog.String("error", err.Error()),
				slog.String("code", st.Code().String()),
			)
		} else {
			logger.Info("gRPC stream completed",
				slog.String("method", info.FullMethod),
				slog.String("duration", duration.String()),
			)
		}

		return err
	}
}
