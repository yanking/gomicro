// Package main HTTP API.
//
// the purpose of this application is to provide an HTTP API for our GoMicro service
//
//	Schemes: http
//	Host: localhost:8080
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
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

// Response represents a generic response.
// swagger:response response
type Response struct {
	// The message
	// in: body
	Message string `json:"message"`
}

// HealthResponse represents a health check response.
// swagger:response healthResponse
type HealthResponse struct {
	// The status
	// in: body
	Status string `json:"status"`
	// The server name
	// in: body
	Server string `json:"server"`
}

// ErrorResponse represents an error response.
// swagger:response ErrorResponse
type ErrorResponse struct {
	// The error message
	// in: body
	Error string `json:"error"`
}

// EchoRequest represents an echo request.
// swagger:parameters echoRequest
type EchoRequest struct {
	// The data to echo
	// in: body
	// required: true
	Data map[string]interface{} `json:"data"`
}

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

	// swagger:route GET / root
	//
	// Returns a welcome message.
	//
	// Responses:
	//   200: response
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, GoMicro HTTP Server!",
		})
	})

	// swagger:route GET /healthz health
	//
	// Returns health status.
	//
	// Responses:
	//   200: healthResponse
	server.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"server": "gomicro-http-example",
		})
	})

	// swagger:route POST /echo echo
	//
	// Echoes the request data.
	//
	// Responses:
	//   200: response
	//   400: ErrorResponse
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
