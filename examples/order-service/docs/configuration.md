# 配置管理

订单服务使用严格的配置加载策略。如果无法加载或解析配置文件，服务将无法启动。这确保了服务始终使用正确的配置运行。

## 命令行参数

服务接受一个命令行参数来指定配置文件路径：

```bash
-config string
    配置文件路径 (默认 "./configs/config.yaml")
```

使用示例：
```bash
# 使用默认配置文件位置
./order-service

# 指定自定义配置文件位置
./order-service -config /path/to/custom/config.yaml
```

## 配置文件结构

配置文件遵循以下结构：

```yaml
server:
  port: "8080"
  host: "localhost"

database:
  - instance: "default"
    driver: "memory"
    host: "localhost"
    port: "3306"
    username: "user"
    password: "password"
    name: "orderdb"
  - instance: "mysql"
    driver: "mysql"
    host: "localhost"
    port: "3306"
    username: "root"
    password: "password"
    name: "orderdb"
  - instance: "mongo"
    driver: "mongo"
    uri: "mongodb://localhost:27017"
    database: "orderdb"
```

### 服务器配置

`server` 部分定义了 HTTP 服务器设置：
- `port`: 服务器监听的端口
- `host`: 服务器绑定的主机地址

### 数据库配置

`database` 部分是数据库配置列表，支持多个实例：
- `instance`: 数据库实例名称（例如，"default", "mysql", "mongo"）
- `driver`: 要使用的数据库驱动（"memory", "mysql", "mongo"）
- `host`: 数据库主机（用于 MySQL）
- `port`: 数据库端口（用于 MySQL）
- `username`: 数据库用户名（用于 MySQL）
- `password`: 数据库密码（用于 MySQL）
- `name`: 数据库名称（用于 MySQL）
- `uri`: 连接 URI（用于 MongoDB）
- `database`: 数据库名称（用于 MongoDB）

## 在代码中使用配置

配置在启动时使用 `config.Load()` 函数加载：

```go
cfg, err := config.Load(*configFile)
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}
```

获取默认数据库配置：

```go
defaultDBConfig := cfg.GetDefaultDatabaseConfig()
```

获取特定数据库实例配置：

```go
mysqlConfig := cfg.GetDatabaseConfig("mysql")
```

## 环境变量

服务还通过 `conf` 包支持环境变量。任何配置值都可以通过设置带有 `GO_KIT_` 前缀并将点号替换为下划线的环境变量来覆盖。

例如：
- `GO_KIT_SERVER_PORT` 覆盖 `server.port`
- `GO_KIT_DATABASE_0_USERNAME` 覆盖第一个数据库实例的用户名

## 错误处理

如果无法加载或解析配置文件，服务将退出并显示错误消息。这种严格的行为确保了服务始终使用有效的配置运行。

## 最佳实践

1. **始终提供配置**: 确保启动服务时始终有有效的配置文件可用。

2. **验证配置**: 检查所有必需的配置值是否存在且有效。

3. **环境特定配置**: 为不同环境（开发、预发布、生产）使用不同的配置文件。

4. **敏感信息**: 永远不要将敏感信息如密码提交到版本控制中。对敏感数据使用环境变量或安全的配置管理系统。

5. **文档**: 保持配置文件的良好文档，以便其他开发人员了解可用选项。