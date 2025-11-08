// Package clientinterceptors provides common gRPC client interceptors.
package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// RetryInterceptor returns a new unary client interceptor that retries failed requests.
func RetryInterceptor(maxRetries int, retryInterval time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var lastErr error
		for i := 0; i <= maxRetries; i++ {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
			lastErr = err
			if i < maxRetries {
				timer := time.NewTimer(retryInterval)
				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				case <-timer.C:
				}
			}
		}
		return lastErr
	}
}
