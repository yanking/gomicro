# asynqctl

`asynqctl` 是一个命令行工具，用于检查和管理 Asynq 队列和任务。

## 安装

确保在项目根目录下运行以下命令来构建工具：

```bash
go build -o asynqctl cmd/asynqctl/*.go
```

## 使用方法

```bash
asynqctl [command]
```

## 可用命令

### queue - 队列检查命令

#### queue ls - 列出所有队列

列出所有队列及其当前统计信息。

```bash
asynqctl queue ls [flags]
```

Flags:
- `-a, --redis-addr string` - Redis 服务器地址 (默认 "127.0.0.1:6379")
- `-d, --redis-db int` - Redis 数据库编号
- `-u, --redis-username string` - Redis 用户名
- `-p, --redis-password string` - Redis 密码

示例：
```bash
# 使用默认 Redis 配置列出所有队列
./asynqctl queue ls

# 使用自定义 Redis 地址
./asynqctl queue ls -a localhost:6379

# 使用 Redis 数据库和认证信息
./asynqctl queue ls -a localhost:6379 -d 1 -u myuser -p mypassword
```

### task - 任务检查命令

#### task active - 列出活跃任务

列出队列中的活跃任务。

```bash
asynqctl task active [queue] [flags]
```

#### task pending - 列出待处理任务

列出队列中的待处理任务。

```bash
asynqctl task pending [queue] [flags]
```

#### task scheduled - 列出计划任务

列出队列中的计划任务。

```bash
asynqctl task scheduled [queue] [flags]
```

#### task retry - 列出重试任务

列出队列中的重试任务。

```bash
asynqctl task retry [queue] [flags]
```

#### task archived - 列出已归档任务

列出队列中的已归档任务。

```bash
asynqctl task archived [queue] [flags]
```

所有任务命令的通用 Flags:
- `-a, --redis-addr string` - Redis 服务器地址 (默认 "127.0.0.1:6379")
- `-d, --redis-db int` - Redis 数据库编号
- `-u, --redis-username string` - Redis 用户名
- `--redis-password string` - Redis 密码
- `-s, --page-size int` - 页面大小 (默认 30)
- `-p, --page int` - 页码 (从 0 开始)

示例：
```bash
# 列出 default 队列中的活跃任务
./asynqctl task active default

# 列出 critical 队列中的待处理任务，每页显示 10 条记录
./asynqctl task pending critical -s 10

# 列出 retry 队列中的重试任务，显示第 2 页
./asynqctl task retry default -s 20 -p 1
```

### scheduler - 调度器检查命令

#### scheduler ls - 列出调度条目

列出所有调度器条目。

```bash
asynqctl scheduler ls [flags]
```

Flags:
- `-a, --redis-addr string` - Redis 服务器地址 (默认 "127.0.0.1:6379")
- `-d, --redis-db int` - Redis 数据库编号
- `-u, --redis-username string` - Redis 用户名
- `--redis-password string` - Redis 密码

示例：
```bash
# 列出所有调度条目
./asynqctl scheduler ls

# 使用自定义 Redis 配置
./asynqctl scheduler ls -a localhost:6379 -u myuser -p mypassword
```

## 全局 Flags

- `-h, --help` - 显示帮助信息

## 使用示例

```bash
# 查看所有可用命令
./asynqctl --help

# 查看队列统计信息
./asynqctl queue ls

# 查看特定队列中的任务
./asynqctl task active default
./asynqctl task pending critical

# 查看调度器条目
./asynqctl scheduler ls
```