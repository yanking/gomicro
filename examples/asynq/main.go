package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/yanking/gomicro/pkg/client/mq"
)

// Task payloads
type WelcomeEmailPayload struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

type ImageProcessingPayload struct {
	ImageURL string `json:"image_url"`
}

func main() {
	// 创建日志记录器
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 初始化 Asynq 客户端 - 单机模式示例
	asynqOpts := &mq.AsynqOptions{
		Instance:    "default",
		RedisAddr:   "127.0.0.1:6379",
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	}

	// 集群模式示例（注释掉上面的单机模式配置，取消注释下面的代码）
	/*
		asynqOpts := &mq.AsynqOptions{
			Instance:   "default",
			RedisAddrs: []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"},
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		}
	*/

	// 哨兵模式示例（注释掉上面的配置，取消注释下面的代码）
	/*
		asynqOpts := &mq.AsynqOptions{
			Instance:   "default",
			RedisAddrs: []string{"127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"},
			MasterName: "mymaster",
			RedisDB:    0,
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		}
	*/

	client, err := mq.InitAsynq(asynqOpts)
	if err != nil {
		log.Fatal("Failed to initialize Asynq client:", err)
	}

	// 创建 quit 通道
	quit := make(chan struct{})

	// 启动任务派发器
	go dispatchTasks(client, logger, quit)

	// 启动任务处理器
	go processTasks(asynqOpts, logger, quit)

	// 等待中断信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// 通知任务派发器停止
	close(quit)

	logger.Info("Shutting down...")
}

func dispatchTasks(client *asynq.Client, logger *slog.Logger, quit <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 派发欢迎邮件任务
			payload := WelcomeEmailPayload{
				UserID: 123,
				Email:  "user@example.com",
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				logger.Error("Failed to marshal payload", "error", err)
				continue
			}

			task := asynq.NewTask("email:welcome", payloadBytes)
			info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.Queue("default"))
			if err != nil {
				logger.Error("Failed to enqueue task", "error", err)
				continue
			}

			logger.Info("Enqueued welcome email task", "task_id", info.ID)

			// 派发图片处理任务
			imagePayload := ImageProcessingPayload{
				ImageURL: "https://example.com/image.jpg",
			}

			imagePayloadBytes, err := json.Marshal(imagePayload)
			if err != nil {
				logger.Error("Failed to marshal image payload", "error", err)
				continue
			}

			imageTask := asynq.NewTask("image:process", imagePayloadBytes)
			info, err = client.Enqueue(imageTask, asynq.MaxRetry(2), asynq.Queue("low"))
			if err != nil {
				logger.Error("Failed to enqueue image task", "error", err)
				continue
			}

			logger.Info("Enqueued image processing task", "task_id", info.ID)
		case <-quit:
			return
		}
	}
}

func processTasks(opts *mq.AsynqOptions, logger *slog.Logger, quit <-chan struct{}) {
	// 初始化 Asynq 服务器
	server, err := mq.InitAsynqServer(opts)
	if err != nil {
		log.Fatal("Failed to initialize Asynq server:", err)
	}

	// 创建任务处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc("email:welcome", handleWelcomeEmail)
	mux.HandleFunc("image:process", handleImageProcessing)

	// 用于等待服务器停止的 WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)

	// 在 goroutine 中运行服务器
	go func() {
		defer wg.Done()
		if err := server.Run(mux); err != nil {
			logger.Error("Asynq server error", "error", err)
		}
	}()

	logger.Info("Asynq server started")

	// 等待中断信号
	<-quit

	// 优雅停止服务器
	server.Stop()
	logger.Info("Asynq server stopped")
}

func handleWelcomeEmail(ctx context.Context, task *asynq.Task) error {
	var payload WelcomeEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Sending welcome email to user_id=%d, email=%s", payload.UserID, payload.Email)

	// 模拟处理时间
	time.Sleep(2 * time.Second)

	log.Printf("Welcome email sent successfully to %s", payload.Email)
	return nil
}

func handleImageProcessing(ctx context.Context, task *asynq.Task) error {
	var payload ImageProcessingPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Processing image: %s", payload.ImageURL)

	// 模拟处理时间
	time.Sleep(5 * time.Second)

	log.Printf("Image processed successfully: %s", payload.ImageURL)
	return nil
}
