# HTTP 传输层示例

这个示例演示了如何使用 GoMicro 的 HTTP 传输层组件创建一个功能完整的 Web 服务器。

## 功能特性

1. 基于 Gin 框架的 HTTP 服务器
2. CORS 跨域支持
3. 请求上下文管理（日志记录、请求追踪）
4. 健康检查端点
5. JSON 数据处理
6. 优雅关闭机制

## 运行示例

```bash
cd /Users/wangyan/Code/Self/gomicro/examples/http
go run main.go
```

服务器将在 `localhost:8080` 启动。

## 测试 API

### 1. 主页
```bash
curl http://localhost:8080/
```

响应:
```json
{
  "message": "Hello, GoMicro HTTP Server!"
}
```

### 2. 健康检查
```bash
curl http://localhost:8080/healthz
```

响应:
```json
{
  "status": "ok",
  "server": "gomicro-http-example"
}
```

### 3. Echo 接口
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, World!", "timestamp": 1234567890}'
```

响应:
```json
{
  "message": "Hello, World!",
  "timestamp": 1234567890
}
```

## 优雅关闭

要测试优雅关闭功能，可以按 `Ctrl+C` 停止服务器。服务器会等待最多5秒钟来完成正在进行的请求，然后关闭。

## 代码说明

### 服务器配置
```go
server := http.NewServer(
    logger,
    http.WithAddr(":8080"),
    http.WithMode(gin.DebugMode),
    http.WithReadTimeout(5*time.Second),
    http.WithWriteTimeout(10*time.Second),
    http.WithIdleTimeout(15*time.Second),
)
```

这里配置了：
- 监听地址: `:8080`
- 运行模式: 调试模式
- 读取超时: 5秒
- 写入超时: 10秒
- 空闲连接超时: 15秒

### 中间件注册
```go
server.Use(middlewares.Cors)
server.Use(middlewares.Context(logger))
```

注册了两个中间件：
1. `CORS` - 处理跨域请求
2. `Context` - 提供请求上下文管理，包括日志记录和请求追踪

### 路由注册
示例中注册了三个路由：
1. `GET /` - 返回欢迎信息
2. `GET /healthz` - 健康检查端点
3. `POST /echo` - 回显接收到的 JSON 数据

### 优雅关闭
通过信号监听实现优雅关闭：
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 5秒超时关闭
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Stop(ctx)
```