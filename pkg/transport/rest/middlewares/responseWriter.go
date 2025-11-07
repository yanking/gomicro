package middlewares

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装gin的ResponseWriter以捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应数据并捕获到缓冲区
func (r *responseWriter) Write(b []byte) (int, error) {
	// 将数据写入缓冲区
	r.body.Write(b)
	// 写入原始响应写入器
	return r.ResponseWriter.Write(b)
}

// WriteHeader 写入响应头
func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}
