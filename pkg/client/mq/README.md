# Asynq Client

Asynq 是一个 Go 语言编写的简单、可靠、高效的分布式任务队列库，使用 Redis 作为消息代理。

## 功能特性

- 简单易用的 API
- 可靠的任务处理
- 支持多种任务类型
- 任务重试机制
- 任务调度功能
- 多实例支持
- 支持 Redis 单机、集群和哨兵模式
- 与项目中现有的 Redis 客户端集成

## 安装

Asynq 已经作为项目依赖添加，无需额外安装。

## 配置

### 单机模式

在配置文件中添加 Asynq 配置：

```yaml
asynq:
  - instance: "default"
    redis_addr: "127.0.0.1:6379"
    redis_db: 0
    redis_username: ""
    redis_password: ""
    concurrency: 10
    queues:
      critical: 6
      default: 3
      low: 1
```

### Redis 集群模式

```yaml
asynq:
  - instance: "default"
    redis_addrs: 
      - "127.0.0.1:7000"
      - "127.0.0.1:7001"
      - "127.0.0.1:7002"
    redis_username: ""
    redis_password: ""
    concurrency: 10
    queues:
      critical: 6
      default: 3
      low: 1
```

### Redis 哨兵模式

```yaml
asynq:
  - instance: "default"
    redis_addrs: 
      - "127.0.0.1:26379"
      - "127.0.0.1:26380"
      - "127.0.0.1:26381"
    master_name: "mymaster"
    redis_db: 0
    redis_username: ""
    redis_password: ""
    sentinel_username: ""
    sentinel_password: ""
    concurrency: 10
    queues:
      critical: 6
      default: 3
      low: 1
```

## 使用方法

### 1. 初始化 Asynq 客户端

```go
import (
    "github.com/yanking/gomicro/pkg/client/mq"
)

// 单机模式
opts := &mq.AsynqOptions{
    Instance:      "default",
    RedisAddr:     "127.0.0.1:6379",
    RedisDB:       0,
    RedisUsername: "",
    RedisPassword: "",
    Concurrency:   10,
    Queues: map[string]int{
        "critical": 6,
        "default":  3,
        "low":      1,
    },
}

// 集群模式
opts := &mq.AsynqOptions{
    Instance:      "default",
    RedisAddrs:    []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"},
    RedisUsername: "",
    RedisPassword: "",
    Concurrency:   10,
    Queues: map[string]int{
        "critical": 6,
        "default":  3,
        "low":      1,
    },
}

// 哨兵模式
opts := &mq.AsynqOptions{
    Instance:         "default",
    RedisAddrs:       []string{"127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"},
    MasterName:       "mymaster",
    RedisDB:          0,
    RedisUsername:    "",
    RedisPassword:    "",
    SentinelUsername: "",
    SentinelPassword: "",
    Concurrency:      10,
    Queues: map[string]int{
        "critical": 6,
        "default":  3,
        "low":      1,
    },
}

// 初始化 Asynq 客户端
client, err := mq.InitAsynq(opts)
if err != nil {
    log.Fatal(err)
}
```

### 2. 创建和派发任务

```go
// 创建任务
task := asynq.NewTask("email:welcome", map[string]interface{}{
    "user_id": 123,
    "email":   "user@example.com",
})

// 派发任务
info, err := client.Enqueue(task)
if err != nil {
    log.Fatal(err)
}
log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
```

### 3. 处理任务

```go
// 创建任务处理器
mux := asynq.NewServeMux()
mux.HandleFunc("email:welcome", func(ctx context.Context, task *asynq.Task) error {
    var payload struct {
        UserID int    `json:"user_id"`
        Email  string `json:"email"`
    }
    
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
    }
    
    log.Printf("Sending welcome email to user_id=%d, email=%s", payload.UserID, payload.Email)
    // 处理业务逻辑
    
    return nil
})

// 初始化 Asynq 服务器
server, err := mq.InitAsynqServer(opts)
if err != nil {
    log.Fatal(err)
}

// 启动服务器
if err := server.Run(mux); err != nil {
    log.Fatal(err)
}
```

### 4. 定时任务调度

```go
// 初始化调度器
schedulerOpts := &mq.AsynqOptions{
    Instance:    "scheduler",
    RedisAddr:   "127.0.0.1:6379",
    SchedulerOpts: &asynq.SchedulerOpts{
        Location: time.Local,
    },
}

scheduler, err := mq.InitScheduler(schedulerOpts)
if err != nil {
    log.Fatal(err)
}

// 创建定时任务
task := asynq.NewTask("report:generate", map[string]interface{}{
    "type": "daily",
    "date": time.Now().Format("2006-01-02"),
})

// 注册定时任务 (使用 Cron 表达式)
// 每天凌晨2点执行
entryID, err := scheduler.Register("0 2 * * *", task)
if err != nil {
    log.Fatal(err)
}

// 启动调度器
if err := scheduler.Run(); err != nil {
    log.Fatal(err)
}
```

## 多实例支持

项目支持多个 Asynq 实例，每个实例可以有不同的配置：

```go
// 初始化多个实例
opts1 := &mq.AsynqOptions{
    Instance:  "email",
    RedisAddr: "127.0.0.1:6379",
    // ... 其他配置
}

opts2 := &mq.AsynqOptions{
    Instance:   "image_processing",
    RedisAddrs: []string{"127.0.0.1:7000", "127.0.0.1:7001"},
    // ... 其他配置
}

mq.InitAsynq(opts1)
mq.InitAsynq(opts2)

// 获取特定实例
emailClient := mq.GetAsynq("email")
imageClient := mq.GetAsynq("image_processing")

// 获取服务器实例
emailServer := mq.GetAsynqServer("email")

// 获取调度器实例
scheduler := mq.GetScheduler("scheduler")
```

## 完整示例

### 普通任务处理示例
查看 [普通任务处理示例代码](../../examples/asynq/main.go) 了解如何使用 Asynq 客户端和服务器。

### 定时任务调度示例
查看 [定时任务调度示例代码](../../examples/asynq_scheduler/main.go) 了解如何使用 Asynq 调度器创建和管理定时任务。

## 最佳实践

1. 为不同类型的任务使用不同的队列
2. 合理设置并发数以避免系统过载
3. 为任务处理函数添加适当的错误处理和重试逻辑
4. 监控任务处理的性能和错误率
5. 使用任务调度功能处理定时任务
6. 根据实际需求选择合适的 Redis 部署模式（单机/集群/哨兵）
7. 为定时任务使用标准的 Cron 表达式格式
8. 合理设置任务的超时时间和重试次数