# HTTP 传输层

基于 Gin 框架实现的 HTTP 服务器组件，提供了丰富的功能和灵活的配置选项。

## 特性

1. 基于 Gin 框架构建
2. 支持中间件机制
3. 优雅启动和关闭
4. 健康检查端点
5. 可配置的超时设置
6. 请求日志记录
7. 跨域支持
8. pprof 性能分析（可选）

## 安装

```bash
go get github.com/yanking/gomicro/pkg/transport/rest
```

## 基本用法

### 创建简单的 HTTP 服务器

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/yanking/gomicro/pkg/transport/rest"
)

func main() {
    // 创建日志记录器
    logger := slog.Default()
    
    // 创建 HTTP 服务器
    server := rest.NewServer(
        logger,
        rest.WithAddr(":8080"),
        rest.WithMode(gin.ReleaseMode),
    )
    
    // 注册路由
    server.GET("/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello, World!",
        })
    })
    
    // 启动服务器
    if err := server.Start(context.Background()); err != nil {
        logger.Error("Failed to start server", "error", err)
    }
}
```

### 使用选项函数配置服务器

```go
server := rest.NewServer(
    logger,
    rest.WithAddr(":9000"),                    // 设置监听地址
    rest.WithMode(gin.ReleaseMode),            // 设置运行模式
    rest.WithReadTimeout(time.Second * 5),     // 设置读取超时
    rest.WithWriteTimeout(time.Second * 10),   // 设置写入超时
    rest.WithIdleTimeout(time.Second * 15),    // 设置空闲连接超时
    rest.WithMaxHeaderBytes(1 << 20),          // 设置最大请求头大小
    rest.WithTrustedProxies([]string{"192.168.1.0/24"}), // 设置受信任代理
    rest.WithHealthz(true),                    // 启用健康检查
    rest.WithEnableProfiling(true),            // 启用性能分析
    rest.WithTransName("zh"),                  // 设置翻译语言
)
```

## 中间件使用

### 使用内置中间件

```go
import "github.com/yanking/gomicro/pkg/transport/rest/middlewares"

server := rest.NewServer(logger)

// 使用 CORS 中间件
server.Use(middlewares.Cors)

// 使用上下文中间件（提供请求追踪和日志记录）
server.Use(middlewares.Context(logger))
```

### 自定义中间件

```go
// 自定义日志中间件
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery
        
        // 处理请求
        c.Next()
        
        // 记录日志
        latency := time.Since(start)
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        
        slog.Info("HTTP request",
            "method", method,
            "path", path,
            "query", raw,
            "ip", clientIP,
            "status", statusCode,
            "latency", latency,
        )
    }
}

// 注册中间件
server.Use(LoggingMiddleware())
```

## 健康检查

默认情况下，服务器会启用 `/healthz` 端点用于健康检查：

```bash
curl http://localhost:8080/healthz
# 返回: {"status":"ok"}
```

可以通过 `WithHealthz(false)` 选项禁用健康检查端点。

## 性能分析

通过 `WithEnableProfiling(true)` 启用 pprof 性能分析，访问以下端点：

- `/debug/pprof/` - pprof 索引页面
- `/debug/pprof/cmdline` - 命令行信息
- `/debug/pprof/profile` - CPU 性能分析
- `/debug/pprof/symbol` - 符号信息
- `/debug/pprof/trace` - trace 信息

## 优雅关闭

服务器支持优雅关闭，确保正在处理的请求能够完成：

```go
// 在单独的 goroutine 中启动服务器
go func() {
    if err := server.Start(context.Background()); err != nil {
        slog.Error("Server error", "error", err)
    }
}()

// 等待中断信号
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 优雅关闭服务器
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := server.Stop(ctx); err != nil {
    slog.Error("Server shutdown error", "error", err)
}
```

## 配置选项

### WithAddr(addr string)
设置服务器监听地址，默认为 ":8080"

### WithMode(mode string)
设置 Gin 运行模式，可选值：
- `gin.DebugMode`
- `gin.ReleaseMode` 
- `gin.TestMode`

### WithReadTimeout(timeout time.Duration)
设置读取超时时间，默认为 5 秒

### WithWriteTimeout(timeout time.Duration)
设置写入超时时间，默认为 10 秒

### WithIdleTimeout(timeout time.Duration)
设置空闲连接超时时间，默认为 15 秒

### WithMaxHeaderBytes(bytes int)
设置最大请求头字节数，默认为 1MB

### WithTrustedProxies(proxies []string)
设置受信任的代理地址列表

### WithEnableProfiling(profiling bool)
启用或禁用 pprof 性能分析，默认为 true

### WithHealthz(healthz bool)
启用或禁用健康检查端点，默认为 true

### WithMetrics(enable bool)
启用或禁用指标收集，默认为 true

### WithTransName(transName string)
设置翻译器语言，默认为 "zh"

## 示例代码

详细示例请参考 [examples/http](../../examples/http) 目录。

## 最佳实践

1. **合理设置超时时间**：根据业务需求设置合适的超时时间
2. **使用中间件**：利用中间件处理跨域、日志记录等通用功能
3. **启用健康检查**：在生产环境中启用健康检查端点
4. **优雅关闭**：确保服务器能够优雅关闭，避免请求中断
5. **日志记录**：使用结构化日志记录关键信息
6. **性能监控**：在需要时启用 pprof 进行性能分析