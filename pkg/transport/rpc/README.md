# RPC 服务支持

## 概述

RPC 包提供了基于 Google gRPC 框架的服务实现，具有以下特性：

- 优雅启动和关闭
- 健康检查支持
- gRPC 反射支持
- TLS 安全支持
- 拦截器支持
- 与项目其他组件一致的 API 设计

## 使用方法

### 1. 创建 gRPC 服务器

```go
// 创建日志记录器
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// 创建gRPC服务器
server := rpc.NewServer(logger, 
    rpc.WithAddress(":9000"),
    rpc.WithHealthz(true),
    rpc.WithReflection(true),
)

// 注册你的gRPC服务
// pb.RegisterYourServiceServer(server.Server, &YourService{})

// 启动服务器
if err := server.Start(context.Background()); err != nil {
    log.Fatal("Failed to start gRPC server:", err)
}
```

### 2. 创建 gRPC 客户端

```go
// 创建日志记录器
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// 创建gRPC客户端
client, err := rpc.NewClient(logger, "localhost:9000",
    rpc.WithTimeout(5*time.Second),
    rpc.WithInsecure(),
)
if err != nil {
    log.Fatal("Failed to create gRPC client:", err)
}
defer client.Close()

// 使用客户端连接调用gRPC服务
// yourClient := pb.NewYourServiceClient(client.GetConn())
```

### 3. 服务器配置选项

```go
// TLS 配置
tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
}
server := rpc.NewServer(logger, rpc.WithTLS(tlsConfig))

// 拦截器
server := rpc.NewServer(logger, 
    rpc.WithUnaryInterceptor(loggingInterceptor),
    rpc.WithStreamInterceptor(authInterceptor),
)

// 多个拦截器
server := rpc.NewServer(logger, 
    rpc.WithUnaryInterceptors(loggingInterceptor, authInterceptor, metricsInterceptor),
)
```

### 4. 客户端配置选项

```go
// TLS 配置
tlsConfig := &tls.Config{
    // TLS配置
}
client, err := rpc.NewClient(logger, "localhost:9000",
    rpc.WithTLS(tlsConfig),
    rpc.WithTimeout(5*time.Second),
)

// 拦截器
client, err := rpc.NewClient(logger, "localhost:9000",
    rpc.WithClientUnaryInterceptor(loggingInterceptor),
    rpc.WithClientStreamInterceptor(metricsInterceptor),
)

// 多个拦截器
client, err := rpc.NewClient(logger, "localhost:9000",
    rpc.WithClientUnaryInterceptors(loggingInterceptor, retryInterceptor, metricsInterceptor),
)
```

### 5. 与应用框架集成

```go
// 创建服务上下文
serviceContext := &ServiceContext{
    Logger: logger,
}

// 创建应用
app, err := app.New(serviceContext, "my-app", "1.0.0", server)
if err != nil {
    log.Fatal("Failed to create app:", err)
}

// 运行应用
if err := app.Run(); err != nil {
    log.Fatal("Failed to run app:", err)
}
```

### 6. 健康检查

```go
// 获取健康检查服务器实例
healthServer := server.GetHealthServer()

// 设置服务状态
healthServer.SetServingStatus("my.service", healthpb.HealthCheckResponse_SERVING)
```

## 服务器拦截器

### 1. 日志记录拦截器

记录请求和响应的日志信息。

### 2. 恢复拦截器

从 panic 中恢复并记录错误信息。

### 3. 超时拦截器

为请求设置超时时间。

### 4. 认证拦截器

验证请求的授权令牌。

### 5. 指标收集拦截器

收集请求的指标信息。

## 客户端拦截器

### 1. 日志记录拦截器

记录客户端请求和响应的日志信息。

### 2. 超时拦截器

为请求设置超时时间。

### 3. 重试拦截器

在请求失败时自动重试。

### 4. 认证拦截器

添加认证令牌到请求中。

### 5. 指标收集拦截器

收集请求的指标信息。

## 最佳实践

1. 在应用程序启动时初始化 gRPC 服务器
2. 使用拦截器实现日志记录、认证、监控等功能
3. 启用健康检查以支持服务发现和负载均衡
4. 在生产环境中使用 TLS 加密通信
5. 合理设置超时和重试策略