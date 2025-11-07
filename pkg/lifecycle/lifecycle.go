// Package lifecycle provides interfaces for component lifecycle management.
// It defines the Component interface that applications can implement to manage
// the lifecycle of their components.
package lifecycle

import "context"

// Component 定义组件的生命周期接口
// 实现该接口的组件可以被 App 框架管理其启动和停止过程
type Component interface {
	// Start 启动组件
	// ctx 是上下文，可用于控制组件的生命周期
	Start(ctx context.Context) error

	// Stop 停止组件
	// ctx 是上下文，可用于控制组件停止的超时
	Stop(ctx context.Context) error

	// Name 返回组件的名称
	// 用于日志记录和调试
	Name() string
}
