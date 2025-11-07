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