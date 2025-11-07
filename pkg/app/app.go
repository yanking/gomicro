// Package app provides application lifecycle management functionality.
// It offers a framework for building applications with component management,
// graceful shutdown, signal handling, and resource cleanup.
package app

import (
	"context"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yanking/gomicro/pkg/lifecycle"
)

const shutdownOverallTimeout = 30 * time.Second

// IConfigProvider 定义配置提供者的接口
type IConfigProvider interface {
	// 这里定义需要访问的配置方法，根据实际需要定义
	// 例如：
	// GetServerPort() int
	// GetDatabaseDSN() string
}

// IServiceContext 定义服务上下文接口
type IServiceContext interface {
	GetLogger() *slog.Logger
	GetConfig() IConfigProvider
}

// Close 定义应用关闭时需要执行的清理函数类型
type Close func(ctx context.Context) error

// App 代表一个应用实例
type App struct {
	appName        string
	serviceVersion string
	logger         *slog.Logger
	cfg            IConfigProvider
	components     []lifecycle.Component
	extCloses      []Close
	runWg          sync.WaitGroup
}

// New 创建一个新的应用实例
func New(ctx IServiceContext, appName, version string, components ...lifecycle.Component) (*App, error) {
	app := &App{
		appName:        appName,
		serviceVersion: version,
		logger:         ctx.GetLogger().With(slog.String("component", "app")),
		cfg:            ctx.GetConfig(),
		components:     components,
	}

	app.logger.Info("Creating new application instance...",
		slog.String("appName", app.appName),
		slog.String("version", app.serviceVersion),
	)
	return app, nil
}

// RegisterClose 注册应用关闭时需要执行的清理函数
func (a *App) RegisterClose(closer ...Close) {
	a.extCloses = append(a.extCloses, closer...)
}

// Run 启动并运行应用
func (a *App) Run() error {
	a.logger.Info("Starting application...")

	// 创建可被信号中断的上下文
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 启动所有组件
	a.startComponents(appCtx)

	a.logger.Info("All components started, application is running.")

	// 等待中断信号
	<-appCtx.Done()

	// 处理关闭流程
	return a.shutdown()
}

// startComponents 启动所有组件
func (a *App) startComponents(ctx context.Context) {
	for _, component := range a.components {
		a.runWg.Add(1)
		go func(component lifecycle.Component) {
			defer a.runWg.Done()
			if err := component.Start(ctx); err != nil {
				a.logger.Error("Component failed to start",
					slog.String("component", component.Name()),
					slog.Any("error", err),
				)
			}
		}(component)
	}
}

// shutdown 优雅关闭应用
func (a *App) shutdown() error {
	a.logger.Info("Shutdown signal received, stopping application.")

	// 创建用于关闭流程的上下文
	stopCtx, cancel := context.WithTimeout(context.Background(), shutdownOverallTimeout)
	defer cancel()

	// 按相反顺序停止组件
	a.stopComponents(stopCtx)

	// 执行外部资源清理
	a.closeExternalResources(stopCtx)

	// 等待所有组件启动goroutine完成
	a.logger.Info("Waiting for component start goroutines to complete...")
	a.runWg.Wait()

	a.logger.Info("Application stopped gracefully.")
	return nil
}

// stopComponents 按相反顺序停止所有组件
func (a *App) stopComponents(ctx context.Context) {
	a.logger.Info("Starting graceful shutdown of components...")

	for i := len(a.components) - 1; i >= 0; i-- {
		component := a.components[i]
		a.logger.Info("Stopping component...", slog.String("component", component.Name()))

		if err := component.Stop(ctx); err != nil {
			a.logger.Error("Component failed to stop",
				slog.String("component", component.Name()),
				slog.Any("error", err),
			)
		} else {
			a.logger.Info("Component stopped successfully",
				slog.String("component", component.Name()),
			)
		}
	}
}

// closeExternalResources 关闭外部资源
func (a *App) closeExternalResources(ctx context.Context) {
	a.logger.Info("Closing external resources...")

	for _, closeFn := range a.extCloses {
		if err := closeFn(ctx); err != nil {
			a.logger.Error("Failed to close resource", slog.Any("error", err))
		}
	}
}
