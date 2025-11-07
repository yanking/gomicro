// Package middlewares provides HTTP middleware functions for the Gin framework.
// It includes context management, logging, CORS support, and response capturing.
package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Context 创建一个上下文中间件，用于日志记录和请求追踪
func Context(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())

		// 创建带请求信息的日志记录器
		ctxLogger := logger.With(
			"request_id", requestID,
			"path", fmt.Sprintf("%s|%s", c.Request.Method, c.Request.URL.Path),
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
