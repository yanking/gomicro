package rpc

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// ServerOption 定义gRPC服务器选项函数
type ServerOption func(*Server)

// WithAddress 设置服务器监听地址
func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithTLS 设置TLS配置
func WithTLS(tlsConfig *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConfig = tlsConfig
	}
}

// WithUnaryInterceptor 添加一元拦截器
func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInterceptors = append(s.unaryInterceptors, interceptor)
	}
}

// WithUnaryInterceptors 添加多个一元拦截器
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptor 添加流拦截器
func WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInterceptors = append(s.streamInterceptors, interceptor)
	}
}

// WithStreamInterceptors 添加多个流拦截器
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInterceptors = append(s.streamInterceptors, interceptors...)
	}
}

// WithHealthz 启用/禁用健康检查服务
func WithHealthz(enabled bool) ServerOption {
	return func(s *Server) {
		s.healthz = enabled
	}
}

// WithReflection 启用/禁用gRPC反射
func WithReflection(enabled bool) ServerOption {
	return func(s *Server) {
		s.enableReflection = enabled
	}
}

// ClientOption 定义gRPC客户端选项函数
type ClientOption func(*Client)

// WithTimeout 设置连接超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithInsecure 使用不安全连接（无TLS）
func WithInsecure() ClientOption {
	return func(c *Client) {
		c.tlsConfig = nil
	}
}

// WithClientUnaryInterceptor 添加一元拦截器
func WithClientUnaryInterceptor(interceptor grpc.UnaryClientInterceptor) ClientOption {
	return func(c *Client) {
		c.unaryInterceptors = append(c.unaryInterceptors, interceptor)
	}
}

// WithClientUnaryInterceptors 添加多个一元拦截器
func WithClientUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) ClientOption {
	return func(c *Client) {
		c.unaryInterceptors = append(c.unaryInterceptors, interceptors...)
	}
}

// WithClientStreamInterceptor 添加流拦截器
func WithClientStreamInterceptor(interceptor grpc.StreamClientInterceptor) ClientOption {
	return func(c *Client) {
		c.streamInterceptors = append(c.streamInterceptors, interceptor)
	}
}

// WithClientStreamInterceptors 添加多个流拦截器
func WithClientStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) ClientOption {
	return func(c *Client) {
		c.streamInterceptors = append(c.streamInterceptors, interceptors...)
	}
}

// WithKeepaliveParams 设置keepalive参数
func WithKeepaliveParams(params keepalive.ClientParameters) ClientOption {
	return func(c *Client) {
		c.keepaliveParams = params
	}
}
