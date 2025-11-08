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
type ReportGenerationPayload struct {
	ReportType string `json:"report_type"`
	Date       string `json:"date"`
}

type DataCleanupPayload struct {
	DaysToKeep int `json:"days_to_keep"`
}

func main() {
	// 创建日志记录器
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 初始化 Asynq 客户端
	asynqOpts := &mq.AsynqOptions{
		Instance:    "scheduler",
		RedisAddr:   "127.0.0.1:6379",
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	}

	// 初始化 Scheduler
	schedulerOpts := &mq.AsynqOptions{
		Instance:  "scheduler",
		RedisAddr: "127.0.0.1:6379",
		SchedulerOpts: &asynq.SchedulerOpts{
			Location: time.Local,
		},
	}

	scheduler, err := mq.InitScheduler(schedulerOpts)
	if err != nil {
		log.Fatal("Failed to initialize scheduler:", err)
	}

	// 注册定时任务
	registerScheduledTasks(scheduler, logger)

	// 启动任务处理器
	go processTasks(asynqOpts, logger)

	// 启动调度器
	go runScheduler(scheduler, logger)

	// 等待中断信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Info("Shutting down...")
}

// registerScheduledTasks 注册定时任务
func registerScheduledTasks(scheduler *asynq.Scheduler, logger *slog.Logger) {
	// 每天凌晨2点生成日报
	payload, _ := json.Marshal(map[string]interface{}{
		"report_type": "daily",
		"date":        time.Now().Format("2006-01-02"),
	})

	dailyReportTask := asynq.NewTask("report:generate", payload)

	entryID, err := scheduler.Register("0 2 * * *", dailyReportTask, asynq.Queue("default"))
	if err != nil {
		logger.Error("Failed to register daily report task", "error", err)
	} else {
		logger.Info("Registered daily report task", "entry_id", entryID)
	}

	// 每周一凌晨3点生成周报
	payload, _ = json.Marshal(map[string]interface{}{
		"report_type": "weekly",
		"date":        time.Now().Format("2006-01-02"),
	})

	weeklyReportTask := asynq.NewTask("report:generate", payload)

	entryID, err = scheduler.Register("0 3 * * 1", weeklyReportTask, asynq.Queue("default"))
	if err != nil {
		logger.Error("Failed to register weekly report task", "error", err)
	} else {
		logger.Info("Registered weekly report task", "entry_id", entryID)
	}

	// 每月1号凌晨4点生成月报
	payload, _ = json.Marshal(map[string]interface{}{
		"report_type": "monthly",
		"date":        time.Now().Format("2006-01-02"),
	})

	monthlyReportTask := asynq.NewTask("report:generate", payload)

	entryID, err = scheduler.Register("0 4 1 * *", monthlyReportTask, asynq.Queue("low"))
	if err != nil {
		logger.Error("Failed to register monthly report task", "error", err)
	} else {
		logger.Info("Registered monthly report task", "entry_id", entryID)
	}

	// 每天凌晨1点清理30天前的数据
	payload, _ = json.Marshal(map[string]interface{}{
		"days_to_keep": 30,
	})

	cleanupTask := asynq.NewTask("data:cleanup", payload)

	entryID, err = scheduler.Register("0 1 * * *", cleanupTask, asynq.Queue("low"))
	if err != nil {
		logger.Error("Failed to register data cleanup task", "error", err)
	} else {
		logger.Info("Registered data cleanup task", "entry_id", entryID)
	}

	// 每5分钟检查系统状态
	systemCheckTask := asynq.NewTask("system:check", []byte("{}"))

	entryID, err = scheduler.Register("*/5 * * * *", systemCheckTask, asynq.Queue("critical"))
	if err != nil {
		logger.Error("Failed to register system check task", "error", err)
	} else {
		logger.Info("Registered system check task", "entry_id", entryID)
	}
}

// runScheduler 启动调度器
func runScheduler(scheduler *asynq.Scheduler, logger *slog.Logger) {
	logger.Info("Starting scheduler")

	// 在 goroutine 中运行调度器
	if err := scheduler.Run(); err != nil {
		logger.Error("Scheduler error", "error", err)
	}
}

// processTasks 处理任务
func processTasks(opts *mq.AsynqOptions, logger *slog.Logger) {
	// 初始化 Asynq 服务器
	server, err := mq.InitAsynqServer(opts)
	if err != nil {
		log.Fatal("Failed to initialize Asynq server:", err)
	}

	// 创建任务处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc("report:generate", handleReportGeneration)
	mux.HandleFunc("data:cleanup", handleDataCleanup)
	mux.HandleFunc("system:check", handleSystemCheck)

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
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// 优雅停止服务器
	server.Stop()
	logger.Info("Asynq server stopped")

	// 等待服务器完全停止
	wg.Wait()
}

// handleReportGeneration 处理报告生成任务
func handleReportGeneration(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		ReportType string `json:"report_type"`
		Date       string `json:"date"`
	}

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Generating %s report for %s", payload.ReportType, payload.Date)

	// 模拟处理时间
	time.Sleep(3 * time.Second)

	log.Printf("%s report generated successfully for %s", payload.ReportType, payload.Date)
	return nil
}

// handleDataCleanup 处理数据清理任务
func handleDataCleanup(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		DaysToKeep int `json:"days_to_keep"`
	}

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("Cleaning up data older than %d days", payload.DaysToKeep)

	// 模拟处理时间
	time.Sleep(5 * time.Second)

	log.Printf("Data cleanup completed, kept data from last %d days", payload.DaysToKeep)
	return nil
}

// handleSystemCheck 处理系统检查任务
func handleSystemCheck(ctx context.Context, task *asynq.Task) error {
	log.Printf("Performing system health check")

	// 模拟处理时间
	time.Sleep(1 * time.Second)

	log.Printf("System health check completed")
	return nil
}
