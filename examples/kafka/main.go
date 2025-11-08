package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yanking/gomicro/pkg/client/mq"
	"github.com/yanking/gomicro/pkg/conf"
)

// Config 定义应用配置结构
type Config struct {
	Kafka []*mq.KafkaOptions `mapstructure:"kafka"`
}

func main() {
	// 设置日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 解析配置
	var cfg Config
	if err := conf.Parse("configs/config.yaml", &cfg); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	// 初始化所有 Kafka 实例
	if err := mq.InitKafkas(cfg.Kafka); err != nil {
		log.Fatal("Failed to initialize Kafka instances:", err)
	}

	// 输出 Kafka 实例信息
	producers := mq.GetKafkaProducerInstances()
	consumers := mq.GetKafkaConsumerInstances()
	logger.Info("Kafka instances initialized", "producers", producers, "consumers", consumers)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动消费者
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeMessages(ctx, logger)
	}()

	// 发送一些测试消息
	sendTestMessages(logger)

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")
	cancel()
	wg.Wait()

	// 关闭所有 Kafka 连接
	if err := mq.CloseKafkaProducer(context.Background()); err != nil {
		logger.Error("Failed to close Kafka producers", "error", err)
	}

	if err := mq.CloseKafkaConsumer(context.Background()); err != nil {
		logger.Error("Failed to close Kafka consumers", "error", err)
	}

	logger.Info("Application stopped")
}

func consumeMessages(ctx context.Context, logger *slog.Logger) {
	// 消费指定主题和分区的消息
	messageChan, errorChan := mq.ConsumeMessages(ctx, "default", "test-topic", 0)

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			logger.Info("Received message",
				"topic", msg.Topic,
				"partition", msg.Partition,
				"offset", msg.Offset,
				"key", string(msg.Key),
				"value", string(msg.Value))
		case err, ok := <-errorChan:
			if !ok {
				return
			}
			logger.Error("Error consuming messages", "error", err)
		case <-ctx.Done():
			logger.Info("Consumer context cancelled")
			return
		}
	}
}

func sendTestMessages(logger *slog.Logger) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		counter := 0
		for range ticker.C {
			counter++
			key := []byte("key-" + string(rune(counter+'0')))
			value := []byte("Hello Kafka! Message #" + string(rune(counter+'0')))

			partition, offset, err := mq.SendMessage("default", "test-topic", key, value)
			if err != nil {
				logger.Error("Failed to send message", "error", err)
			} else {
				logger.Info("Message sent",
					"partition", partition,
					"offset", offset,
					"key", string(key),
					"value", string(value))
			}
		}
	}()
}
