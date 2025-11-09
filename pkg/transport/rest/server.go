// Package rest provides an HTTP server implementation based on Gin framework.
// It includes features like middleware support, graceful shutdown, health checks,
// and integration with various tools.
package rest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

// Server 代表HTTP服务器
type Server struct {
	*gin.Engine
	addr           string
	mode           string
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
	maxHeaderBytes int
	trustedProxies []string
	server         *http.Server

	healthz         bool
	enableProfiling bool
	enableMetrics   bool

	transName string
	trans     ut.Translator

	logger *slog.Logger
}

// NewServer 创建一个新的HTTP服务器实例
func NewServer(logger *slog.Logger, opts ...ServerOption) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	srv := &Server{
		Engine:          gin.New(), // 使用New而不是Default，提供更灵活的中间件控制
		addr:            ":8080",
		mode:            gin.DebugMode,
		readTimeout:     5 * time.Second,
		writeTimeout:    10 * time.Second,
		idleTimeout:     15 * time.Second,
		maxHeaderBytes:  1 << 20, // 1 MB
		trustedProxies:  nil,
		healthz:         true,
		enableProfiling: true,
		enableMetrics:   true,
		logger:          logger,
	}

	// 应用选项
	for _, opt := range opts {
		opt(srv)
	}

	// 注册默认中间件
	srv.Engine.Use(gin.Logger())
	srv.Engine.Use(gin.Recovery())

	// 注册健康检查路由
	if srv.healthz {
		srv.Engine.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}

	return srv
}

// Start 启动HTTP服务器
func (s *Server) Start(_ context.Context) error {
	if s.mode != gin.DebugMode && s.mode != gin.ReleaseMode && s.mode != gin.TestMode {
		return errors.New("mode must be one of 'debug', 'release', or 'test'")
	}

	gin.SetMode(s.mode)

	// 初始化翻译器
	if err := s.initTrans(s.transName); err != nil {
		s.logger.Error("Translator init failed", slog.Any("error", err))
		return err
	}

	s.logger.Info("Starting HTTP server", slog.String("addr", s.addr))

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:           s.addr,
		Handler:        s.Engine,
		ReadTimeout:    s.readTimeout,
		WriteTimeout:   s.writeTimeout,
		IdleTimeout:    s.idleTimeout,
		MaxHeaderBytes: s.maxHeaderBytes,
	}

	// 设置受信任代理
	if err := s.Engine.SetTrustedProxies(s.trustedProxies); err != nil {
		s.logger.Error("Failed to set trusted proxies", slog.Any("error", err))
		return err
	}

	// 启动服务器
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("HTTP server error", slog.Any("error", err))
		return err
	}

	return nil
}

// Stop 停止HTTP服务器
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")

	if s.server != nil {
		// 创建一个带超时的上下文，确保服务器能及时关闭
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("HTTP server shutdown error", slog.Any("error", err))
			// 如果优雅关闭失败，强制关闭
			if forceErr := s.server.Close(); forceErr != nil {
				s.logger.Error("HTTP server force close error", slog.Any("error", forceErr))
				return forceErr
			}
			return err
		}
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

// Name 返回组件名称
func (s *Server) Name() string {
	return "gin_rest_server"
}

// initTrans 初始化翻译器
func (s *Server) initTrans(locale string) error {
	// 创建本地化翻译器
	zhT := zh.New()
	enT := en.New()
	uni := ut.New(enT, zhT, enT)

	var ok bool
	s.trans, ok = uni.GetTranslator(locale)
	if !ok {
		return errors.New("translator not found for locale: " + locale)
	}

	s.logger.Info("Translator initialized", slog.String("locale", locale))
	return nil
}
