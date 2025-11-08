// Package serverinterceptors provides common gRPC server interceptors.
package serverinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TimeoutInterceptor returns a new unary server interceptor that sets a timeout for requests.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan struct{})
		var resp interface{}
		var err error

		go func() {
			resp, err = handler(ctxWithTimeout, req)
			close(done)
		}()

		select {
		case <-done:
			return resp, err
		case <-ctxWithTimeout.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "request timeout")
		}
	}
}
