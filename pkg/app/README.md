# App 包

应用程序框架包，提供通用的应用程序生命周期管理功能。

## 特性

1. 应用程序启动和停止管理
2. 组件生命周期管理
3. 优雅关闭机制
4. 信号处理（SIGINT, SIGTERM）
5. 可扩展的清理函数注册

## 设计理念

### 依赖倒置原则

通过接口抽象减少具体实现的依赖，接口定义在 [app.go](file:///Users/wangyan/Code/Self/gomicro/pkg/app/app.go) 文件中：

```go
// 定义服务上下文接口而不是依赖具体实现
type IServiceContext interface {
    GetLogger() *slog.Logger
    GetConfig() IConfigProvider
}
```

### 组件化架构

支持通过组件接口管理各种服务组件：

```go
type Component interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Name() string
}
```

## 使用示例

详细示例请参考 [examples/app](../../examples/app) 目录。

基本用法：

```go
// 创建服务上下文
ctx := &ExampleServiceContext{
    Logger: slog.Default(),
    Config: &ExampleConfig{
        ServerPort: 8080,
        DatabaseDSN: "user:pass@tcp(localhost:3306)/dbname",
    },
}

// 创建应用实例
app, err := app.New(ctx, "myapp", "v1.0.0", httpComponent, databaseComponent)
if err != nil {
    log.Fatal(err)
}

// 注册清理函数
app.RegisterClose(func(ctx context.Context) error {
    // 清理资源
    return nil
})

// 运行应用
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

## 最佳实践

1. 通过接口定义依赖而不是具体实现
2. 组件应实现生命周期接口
3. 合理使用日志记录关键信息
4. 注册必要的清理函数确保资源释放
5. 遵循优雅关闭原则，确保数据一致性