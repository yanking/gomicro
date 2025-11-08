// Package serverinterceptors provides common gRPC server interceptors.
package serverinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// MetricsInterceptor returns a new unary server interceptor that collects metrics.
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// 这里可以集成监控指标收集逻辑
		// 例如：记录请求计数、延迟等指标
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		// 在实际应用中，这里会将指标发送到监控系统
		// 例如：prometheus、statsd等
		_ = duration // 避免未使用变量警告

		return resp, err
	}
}
