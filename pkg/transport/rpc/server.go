// Package rpc provides a gRPC server implementation.
// It includes features like graceful shutdown, health checks,
// and integration with various tools.
package rpc

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server 代表gRPC服务器
type Server struct {
	*grpc.Server
	addr               string
	healthz            bool
	enableReflection   bool
	tlsConfig          *tls.Config
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	logger             *slog.Logger
	healthServer       *health.Server
}

// NewServer 创建一个新的gRPC服务器实例
func NewServer(logger *slog.Logger, opts ...ServerOption) *Server {
	srv := &Server{
		addr:             ":9000",
		healthz:          true,
		enableReflection: true,
		logger:           logger,
		healthServer:     health.NewServer(),
	}

	// 应用选项
	for _, opt := range opts {
		opt(srv)
	}

	// 创建gRPC服务器选项
	serverOpts := []grpc.ServerOption{}

	// 添加TLS配置（如果提供）
	if srv.tlsConfig != nil {
		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(srv.tlsConfig)))
	}

	// 添加一元拦截器链
	if len(srv.unaryInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(srv.unaryInterceptors...))
	}

	// 添加流拦截器链
	if len(srv.streamInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.ChainStreamInterceptor(srv.streamInterceptors...))
	}

	// 创建gRPC服务器
	srv.Server = grpc.NewServer(serverOpts...)

	// 注册健康检查服务
	if srv.healthz {
		healthpb.RegisterHealthServer(srv.Server, srv.healthServer)
		srv.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	}

	// 注册反射服务
	if srv.enableReflection {
		reflection.Register(srv.Server)
	}

	return srv
}

// Start 启动gRPC服务器
func (s *Server) Start(_ context.Context) error {
	s.logger.Info("Starting gRPC server", slog.String("addr", s.addr))

	// 监听TCP端口
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Error("Failed to listen", slog.Any("error", err))
		return err
	}

	// 启动服务器
	if err := s.Server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		s.logger.Error("gRPC server error", slog.Any("error", err))
		return err
	}

	return nil
}

// Stop 停止gRPC服务器
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping gRPC server")

	// 通知健康检查服务服务器正在关闭
	if s.healthServer != nil {
		s.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}

	// 创建一个带超时的上下文，确保服务器能及时关闭
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 优雅地停止服务器
	done := make(chan struct{})
	go func() {
		s.Server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		s.logger.Warn("gRPC server shutdown timeout, forcing stop")
		s.Server.Stop()
	}

	return nil
}

// Name 返回组件名称
func (s *Server) Name() string {
	return "grpc_server"
}

// GetHealthServer 返回健康检查服务器实例
func (s *Server) GetHealthServer() *health.Server {
	return s.healthServer
}
