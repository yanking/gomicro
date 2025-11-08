# Kafka 消息队列支持

## 配置说明

在 `config.yaml` 中配置多个 Kafka 实例：

```yaml
kafka:
  - instance: "default"
    brokers:
      - "localhost:9092"
    version: "2.1.0"
    producer:
      maxRetries: 3
      retryBackoff: "100ms"
      requiredAcks: 1
      timeout: "10s"
    consumer:
      groupID: "my-consumer-group"
      offsetInitial: -1  # -1 for newest, -2 for oldest
      timeout: "250ms"
```

## 使用方法

### 1. 初始化所有实例

```go
// 定义配置结构体
type Config struct {
    Kafka []*mq.KafkaOptions `mapstructure:"kafka"`
}

// 解析配置
var cfg Config
if err := conf.Parse("configs/config.yaml", &cfg); err != nil {
    log.Fatal(err)
}

// 初始化所有Kafka实例
if err := mq.InitKafkas(cfg.Kafka); err != nil {
    log.Fatal(err)
}
```

### 2. 发送消息

```go
// 发送消息到指定主题
partition, offset, err := mq.SendMessage("default", "my-topic", []byte("key"), []byte("Hello Kafka!"))
if err != nil {
    log.Fatal("Failed to send message:", err)
}
log.Printf("Message sent to partition %d at offset %d", partition, offset)
```

### 3. 消费消息

```go
// 创建上下文
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 消费指定主题和分区的消息
messageChan, errorChan := mq.ConsumeMessages(ctx, "default", "my-topic", 0)

// 处理消息
go func() {
    for {
        select {
        case msg := <-messageChan:
            log.Printf("Received message: %s", string(msg.Value))
            // 处理消息
        case err := <-errorChan:
            log.Printf("Error consuming messages: %v", err)
        case <-ctx.Done():
            return
        }
    }
}()
```

### 4. 关闭连接

```go
// 关闭所有生产者实例
mq.CloseKafkaProducer(context.Background())

// 关闭所有消费者实例
mq.CloseKafkaConsumer(context.Background())

// 关闭指定实例
mq.CloseKafkaProducer(context.Background(), "default")
mq.CloseKafkaConsumer(context.Background(), "default")
```

## 最佳实践

1. 在应用程序启动时初始化所有需要的 Kafka 实例
2. 根据业务需求获取对应实例进行消息发送或消费
3. 在程序退出时关闭所有 Kafka 连接
4. 使用上下文控制消费者生命周期
5. 合理处理消费过程中的错误