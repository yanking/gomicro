// Package clientinterceptors provides common gRPC client interceptors.
package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// MetricsInterceptor returns a new unary client interceptor that collects metrics.
func MetricsInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 这里可以集成监控指标收集逻辑
		// 例如：记录请求计数、延迟等指标
		startTime := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(startTime)
		// 在实际应用中，这里会将指标发送到监控系统
		// 例如：prometheus、statsd等
		_ = duration // 避免未使用变量警告

		return err
	}
}
