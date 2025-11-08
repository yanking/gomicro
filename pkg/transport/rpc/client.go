// Package rpc provides a gRPC client implementation.
// It includes features like connection management, load balancing,
// and integration with various tools.
package rpc

import (
	"context"
	"crypto/tls"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// Client 代表gRPC客户端
type Client struct {
	conn               *grpc.ClientConn
	target             string
	timeout            time.Duration
	tlsConfig          *tls.Config
	unaryInterceptors  []grpc.UnaryClientInterceptor
	streamInterceptors []grpc.StreamClientInterceptor
	keepaliveParams    keepalive.ClientParameters
	logger             *slog.Logger
}

// NewClient 创建一个新的gRPC客户端实例
func NewClient(logger *slog.Logger, target string, opts ...ClientOption) (*Client, error) {
	client := &Client{
		target:  target,
		timeout: 10 * time.Second,
		logger:  logger,
		keepaliveParams: keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(client)
	}

	// 创建gRPC客户端选项
	dialOpts := []grpc.DialOption{}

	// 添加TLS配置（如果提供）
	if client.tlsConfig != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(client.tlsConfig)))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}

	// 添加一元拦截器链
	if len(client.unaryInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(client.unaryInterceptors...))
	}

	// 添加流拦截器链
	if len(client.streamInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainStreamInterceptor(client.streamInterceptors...))
	}

	// 添加keepalive参数
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(client.keepaliveParams))

	// 建立连接
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, client.target, dialOpts...)
	if err != nil {
		logger.Error("Failed to dial gRPC server", slog.Any("error", err))
		return nil, err
	}

	client.conn = conn
	return client, nil
}

// GetConn 返回底层的gRPC连接
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭gRPC客户端连接
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// WithTimeout 设置请求超时时间
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

// Timeout 返回当前超时设置
func (c *Client) Timeout() time.Duration {
	return c.timeout
}
