package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yanking/gomicro/pkg/transport/rest"
	"github.com/yanking/gomicro/pkg/transport/rest/middlewares"
)

func main() {
	// 创建日志记录器
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 创建 HTTP 服务器
	server := rest.NewServer(
		logger,
		rest.WithAddr(":8080"),
		rest.WithMode(gin.DebugMode),
		rest.WithReadTimeout(5*time.Second),
		rest.WithWriteTimeout(10*time.Second),
		rest.WithIdleTimeout(15*time.Second),
	)

	// 注册中间件
	server.Use(middlewares.Cors)
	server.Use(middlewares.Context(logger))

	// 注册路由
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, GoMicro HTTP Server!",
		})
	})

	server.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"server": "gomicro-http-example",
		})
	})

	server.POST("/echo", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	// 启动服务器
	go func() {
		logger.Info("Starting HTTP server on :8080")
		if err := server.Start(context.Background()); err != nil {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭服务器
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	logger.Info("Server exited")
}
