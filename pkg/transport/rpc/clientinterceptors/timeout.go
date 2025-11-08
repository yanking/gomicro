// Package clientinterceptors provides common gRPC client interceptors.
package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutInterceptor returns a new unary client interceptor that sets a timeout for requests.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 创建带超时的上下文
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(ctxWithTimeout, method, req, reply, cc, opts...)
	}
}
