package main

import (
	"log/slog"
	"os"
)

func main() {
	// 创建示例配置
	config := &ExampleConfig{
		ServerPort:  8080,
		DatabaseDSN: "user:pass@tcp(localhost:3306)/dbname",
	}

	// 创建服务上下文
	ctx := &ExampleServiceContext{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Config: config,
	}

	// 输出示例信息
	ctx.GetLogger().Info("Example service context created",
		"server_port", ctx.GetConfig().(*ExampleConfig).GetServerPort(),
		"database_dsn", ctx.GetConfig().(*ExampleConfig).GetDatabaseDSN(),
	)
}
