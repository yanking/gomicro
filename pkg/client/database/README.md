# MySQL 多实例支持

## 配置说明

在 `config.yaml` 中配置多个 MySQL 实例：

```yaml
mysql:
  - instance: "default"
    addr: "localhost:3306"
    username: "root"
    password: "root"
    database: "aidoggy"
    maxIdleConnections: 10
    maxOpenConnections: 100
    maxConnectionLifeTime: "30m"
  - instance: "analytics"
    addr: "localhost:3306"
    username: "root"
    password: "root"
    database: "analytics_db"
    maxIdleConnections: 5
    maxOpenConnections: 50
    maxConnectionLifeTime: "15m"
```

## 使用方法

### 1. 初始化所有实例

```go
// 定义配置结构体
type Config struct {
    MySQL []*database.MySQLOptions `mapstructure:"mysql"`
}

// 解析配置
var cfg Config
if err := conf.Parse("configs/config.yaml", &cfg); err != nil {
    log.Fatal(err)
}

// 初始化所有MySQL实例
if err := database.InitMySQLs(cfg.MySQL); err != nil {
    log.Fatal(err)
}
```

### 2. 获取数据库实例

```go
// 获取默认实例
db := database.GetMySQL()

// 获取指定实例
analyticsDB := database.GetMySQL("analytics")
```

### 3. 关闭连接

```go
// 关闭所有实例
database.CloseMySQL(context.Background())

// 关闭指定实例
database.CloseMySQL(context.Background(), "analytics")
```

## 最佳实践

1. 在应用程序启动时初始化所有需要的实例
2. 根据业务需求获取对应实例进行操作
3. 在程序退出时关闭所有数据库连接

# Redis 多实例支持

## 配置说明

在 `config.yaml` 中配置多个 Redis 实例：

```yaml
redis:
  - instance: "default"
    addrs:
      - "localhost:6379"
    db: 0
    username: ""
    password: ""
    poolSize: 10
    minIdleConns: 5
  - instance: "cache"
    addrs:
      - "localhost:6380"
    db: 1
    username: ""
    password: ""
    poolSize: 20
    minIdleConns: 10
```

## 使用方法

### 1. 初始化所有实例

```go
// 定义配置结构体
type Config struct {
    Redis []*database.RedisOptions `mapstructure:"redis"`
}

// 解析配置
var cfg Config
if err := conf.Parse("configs/config.yaml", &cfg); err != nil {
    log.Fatal(err)
}

// 初始化所有Redis实例
if err := database.InitRedises(cfg.Redis); err != nil {
    log.Fatal(err)
}
```

### 2. 获取Redis实例

```go
// 获取默认实例
rdb := database.GetRedis()

// 获取指定实例
cacheRdb := database.GetRedis("cache")
```

### 3. 关闭连接

```go
// 关闭所有实例
database.CloseRedis(context.Background())

// 关闭指定实例
database.CloseRedis(context.Background(), "cache")
```

## 最佳实践

1. 在应用程序启动时初始化所有需要的实例
2. 根据业务需求获取对应实例进行操作
3. 在程序退出时关闭所有Redis连接

# MongoDB 多实例支持

## 配置说明

在 `config.yaml` 中配置多个 MongoDB 实例：

```yaml
mongodb:
  - instance: "default"
    uri: "mongodb://localhost:27017"
    connectTimeout: "10s"
    maxPoolSize: 10
  - instance: "analytics"
    uri: "mongodb://localhost:27018"
    connectTimeout: "10s"
    maxPoolSize: 5
```

## 使用方法

### 1. 初始化所有实例

```go
// 定义配置结构体
type Config struct {
    MongoDB []*database.MongoDBOptions `mapstructure:"mongodb"`
}

// 解析配置
var cfg Config
if err := conf.Parse("configs/config.yaml", &cfg); err != nil {
    log.Fatal(err)
}

// 初始化所有MongoDB实例
if err := database.InitMongoDBs(cfg.MongoDB); err != nil {
    log.Fatal(err)
}
```

### 2. 获取MongoDB实例

```go
// 获取默认实例
client := database.GetMongoDB()

// 获取指定实例
analyticsClient := database.GetMongoDB("analytics")
```

### 3. 关闭连接

```go
// 关闭所有实例
database.CloseMongoDB(context.Background())

// 关闭指定实例
database.CloseMongoDB(context.Background(), "analytics")
```

## 最佳实践

1. 在应用程序启动时初始化所有需要的实例
2. 根据业务需求获取对应实例进行操作
3. 在程序退出时关闭所有MongoDB连接