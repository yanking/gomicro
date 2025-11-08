// Package clientinterceptors provides common gRPC client interceptors.
package clientinterceptors

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor returns a new unary client interceptor that logs outgoing requests.
func LoggingInterceptor(logger *slog.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()

		logger.Info("gRPC client request started",
			slog.String("method", method),
			slog.String("target", cc.Target()),
			slog.String("start_time", startTime.Format(time.RFC3339)),
		)

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC client request failed",
				slog.String("method", method),
				slog.String("target", cc.Target()),
				slog.String("duration", duration.String()),
				slog.String("error", err.Error()),
				slog.String("code", st.Code().String()),
			)
		} else {
			logger.Info("gRPC client request completed",
				slog.String("method", method),
				slog.String("target", cc.Target()),
				slog.String("duration", duration.String()),
			)
		}

		return err
	}
}

// LoggingStreamInterceptor returns a new stream client interceptor that logs outgoing streams.
func LoggingStreamInterceptor(logger *slog.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		startTime := time.Now()

		logger.Info("gRPC client stream started",
			slog.String("method", method),
			slog.String("target", cc.Target()),
			slog.String("start_time", startTime.Format(time.RFC3339)),
		)

		stream, err := streamer(ctx, desc, cc, method, opts...)

		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC client stream failed",
				slog.String("method", method),
				slog.String("target", cc.Target()),
				slog.String("duration", duration.String()),
				slog.String("error", err.Error()),
				slog.String("code", st.Code().String()),
			)
		} else {
			logger.Info("gRPC client stream created",
				slog.String("method", method),
				slog.String("target", cc.Target()),
				slog.String("duration", duration.String()),
			)
		}

		return stream, err
	}
}
