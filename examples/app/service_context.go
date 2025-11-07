package main

import (
	"log/slog"

	"github.com/yanking/gomicro/pkg/app"
)

// ExampleConfig 示例配置结构
type ExampleConfig struct {
	ServerPort  int
	DatabaseDSN string
}

// GetServerPort 获取服务端口
func (c *ExampleConfig) GetServerPort() int {
	return c.ServerPort
}

// GetDatabaseDSN 获取数据库连接字符串
func (c *ExampleConfig) GetDatabaseDSN() string {
	return c.DatabaseDSN
}

// ExampleServiceContext 示例服务上下文
type ExampleServiceContext struct {
	Logger *slog.Logger
	Config *ExampleConfig
}

// GetLogger 获取日志记录器
func (ctx *ExampleServiceContext) GetLogger() *slog.Logger {
	return ctx.Logger
}

// GetConfig 获取配置提供者
func (ctx *ExampleServiceContext) GetConfig() app.IConfigProvider {
	return ctx.Config
}
