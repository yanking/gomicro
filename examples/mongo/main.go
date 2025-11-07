package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yanking/gomicro/pkg/client/database"
	"github.com/yanking/gomicro/pkg/conf"
)

// Config 定义应用配置结构
type Config struct {
	MongoDB []*database.MongoDBOptions `mapstructure:"mongodb"`
}

func main() {
	// 设置日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 解析配置
	var cfg Config
	if err := conf.Parse("configs/config.yaml", &cfg); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	// 初始化所有 MongoDB 实例
	if err := database.InitMongoDBs(cfg.MongoDB); err != nil {
		log.Fatal("Failed to initialize MongoDB instances:", err)
	}

	// 输出 MongoDB 实例信息
	instances := database.GetMongoDBInstances()
	logger.Info("MongoDB instances initialized", "instances", instances)

	// 获取默认 MongoDB 实例
	client := database.GetMongoDB()
	if client == nil {
		log.Fatal("Failed to get default MongoDB instance")
	}

	// 等待中断信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 模拟一些工作
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				logger.Info("Application is running...")
			case <-ctx.Done():
				return
			}
		}
	}()

	logger.Info("MongoDB example application started. Press Ctrl+C to exit.")

	// 等待中断信号
	<-ctx.Done()

	logger.Info("Shutting down...")

	// 关闭所有 MongoDB 连接
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.CloseMongoDB(shutdownCtx); err != nil {
		logger.Error("Failed to close MongoDB connections", "error", err)
	} else {
		logger.Info("MongoDB connections closed successfully")
	}
}
