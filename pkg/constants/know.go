package constants

// 定义 RequestID 相关常量
const (
	// RequestIDKey 是请求ID在上下文和HTTP头中的键名
	RequestIDKey = "request_id"
)

// RequestIDCtx 请求ID上下文
type RequestIDCtx struct {
}
