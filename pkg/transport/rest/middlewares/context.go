// Package middlewares provides HTTP middleware functions for the Gin framework.
// It includes context management, logging, CORS support, and response capturing.
package middlewares

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yanking/gomicro/pkg/constants"
)

// RequestContextKey 定义request_id在context中的键类型
type RequestContextKey string

const (
	// RequestIDKey 是request_id在context中的键
	RequestIDKey RequestContextKey = "request_id"
)

// Context 创建一个上下文中间件，用于日志记录和请求追踪
func Context(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求中是否已存在request_id，如果不存在则生成新的UUID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// 将request_id放入context中，便于后续处理使用
		ctx := context.WithValue(c.Request.Context(), constants.RequestIDCtx{}, requestID)
		c.Request = c.Request.WithContext(ctx)

		// 创建带请求信息的日志记录器
		ctxLogger := logger.With(
			constants.RequestIDKey, requestID,
			"path", c.Request.Method+"|"+c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
		)

		// 读取并记录请求体（如果有的话）
		var reqBody string
		if c.Request.Body != nil {
			reqBytes, _ := io.ReadAll(c.Request.Body)
			if len(reqBytes) > 0 {
				reqBody = string(reqBytes)
				// 恢复请求体，以便后续处理可以读取
				c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
			}
		}

		// 记录请求信息
		ctxLogger.Info("Request received",
			slog.String("body", reqBody),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		)

		// 替换响应写入器以捕获响应体
		c.Writer = &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}

		// 记录处理时间
		startTime := time.Now()
		c.Next()
		latency := time.Since(startTime)

		// 记录响应信息
		errors := make([]string, len(c.Errors))
		for i, err := range c.Errors {
			errors[i] = err.Error()
		}

		ctxLogger.Info("Response sent",
			slog.String("body", c.Writer.(*responseWriter).body.String()),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
			slog.Any("errors", errors),
		)
	}
}
