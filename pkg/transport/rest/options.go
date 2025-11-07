package rest

import (
	"time"
)

// ServerOption 定义HTTP服务器选项函数
type ServerOption func(*Server)

// WithAddr 设置服务器监听地址
func WithAddr(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithMode 设置Gin运行模式
func WithMode(mode string) ServerOption {
	return func(s *Server) {
		s.mode = mode
	}
}

// WithReadTimeout 设置读取超时时间
func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WithWriteTimeout 设置写入超时时间
func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// WithIdleTimeout 设置空闲连接超时时间
func WithIdleTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.idleTimeout = timeout
	}
}

// WithMaxHeaderBytes 设置最大请求头字节数
func WithMaxHeaderBytes(bytes int) ServerOption {
	return func(s *Server) {
		s.maxHeaderBytes = bytes
	}
}

// WithTrustedProxies 设置受信任的代理
func WithTrustedProxies(proxies []string) ServerOption {
	return func(s *Server) {
		s.trustedProxies = proxies
	}
}

// WithEnableProfiling 启用/禁用pprof性能分析
func WithEnableProfiling(profiling bool) ServerOption {
	return func(s *Server) {
		s.enableProfiling = profiling
	}
}

// WithHealthz 启用/禁用健康检查端点
func WithHealthz(healthz bool) ServerOption {
	return func(s *Server) {
		s.healthz = healthz
	}
}

// WithMetrics 启用/禁用指标收集
func WithMetrics(enable bool) ServerOption {
	return func(s *Server) {
		s.enableMetrics = enable
	}
}

// WithTransName 设置翻译器语言
func WithTransName(transName string) ServerOption {
	return func(s *Server) {
		s.transName = transName
	}
}
