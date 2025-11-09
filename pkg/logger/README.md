# 日志包

日志包提供了对 Go 的 `slog` 包的简单封装，用于在整个应用程序中保持一致的日志记录。

## 功能特性

- 简单的初始化和配置
- 支持不同的日志格式（文本、JSON）
- 支持不同的日志级别
- 全局日志记录器访问
- 线程安全操作
- 默认使用 RFC3339 时间格式
- 可选的源文件和行号信息
- 可选的源文件路径基础路径修剪
- 自动基础路径检测

## 使用方法

### 基本用法

```go
import "github.com/yanking/gomicro/pkg/logger"

// 使用默认配置初始化日志记录器
logger.Init(nil)

// 直接使用 slog 方法
log := logger.Get()
log.Info("应用程序已启动")
log.Error("发生错误", "error", err)
```

### 自定义配置

```go
import "github.com/yanking/gomicro/pkg/logger"

config := &logger.Config{
    Level:              slog.LevelDebug,
    Format:             "json",
    Output:             os.Stdout,
    AddSource:          true,
    BasePath:           "/path/to/project/root/",
    AutoDetectBasePath: true,
}

logger.Init(config)
log := logger.Get()
log.Debug("调试信息")
```

### 直接创建日志记录器实例

```go
import "github.com/yanking/gomicro/pkg/logger"

config := &logger.Config{
    Level:  slog.LevelInfo,
    Format: "text",
}

log := logger.New(config)
log.Info("Hello, world!")
```

## API

### 函数

- `Init(config *Config)` - 初始化全局日志记录器
- `Get() *slog.Logger` - 获取全局日志记录器实例
- `DefaultConfig() *Config` - 获取默认配置

### 类型

- `Config` - 日志记录器配置

### 配置字段

- `Level` - 日志级别 (slog.Level)
- `Format` - 日志格式 ("text" 或 "json")
- `Output` - 输出写入器 (默认: os.Stdout)
- `AddSource` - 是否添加源文件和行号
- `BasePath` - 从源文件路径中修剪的基础路径
- `AutoDetectBasePath` - 通过查找 go.mod 文件自动检测基础路径