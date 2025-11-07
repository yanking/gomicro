# 配置解析包 (conf)

该包提供了配置文件解析功能，支持 YAML 格式的配置文件，并能自动绑定环境变量。

## 功能特性

1. 支持 YAML 格式配置文件解析
2. 自动绑定环境变量
3. 支持配置热重载（文件变更时自动重新加载）
4. 线程安全

## 使用方法

### 基本用法

```go
type Config struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        DSN string `mapstructure:"dsn"`
    } `mapstructure:"database"`
}

var cfg Config
if err := conf.Parse("config.yaml", &cfg); err != nil {
    log.Fatal(err)
}

// 访问解析后的值
fmt.Println("Server:", cfg.Server.Host, cfg.Server.Port)
fmt.Println("Database:", cfg.Database.DSN)
```

### 环境变量绑定

该包会自动绑定环境变量，默认使用 "GO_KIT" 作为前缀，并将配置中的点号(.)替换为下划线(_)。

例如，对于上面的结构体：
- `GO_KIT_SERVER_HOST` 会绑定到 `cfg.Server.Host`
- `GO_KIT_DATABASE_DSN` 会绑定到 `cfg.Database.DSN`

### 配置热重载

支持配置文件变更时自动重新加载：

```go
func reloadCallback() {
    // 处理配置变更
    fmt.Println("配置已重新加载")
}

if err := conf.Parse("config.yaml", &cfg, reloadCallback); err != nil {
    log.Fatal(err)
}
```

在上述例子中，每当 `config.yaml` 文件被修改时，`reloadCallback` 函数都会被执行。

## 注意事项

1. 传入的 `obj` 必须是指针类型
2. 结构体字段需要使用 `mapstructure` 标签来映射配置文件中的键名
3. 热重载功能是可选的，只有在提供了 reload 函数时才会启用