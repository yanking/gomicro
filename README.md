# GoMicro

一个支持多 MySQL 实例的 Go 项目。

## 功能特性

- 支持多 MySQL 实例
- 支持多 Redis 实例
- 配置管理
- 使用 golangci-lint 进行代码检查
- Git pre-commit 钩子
- 基于 Gin 的 HTTP 传输层
- 使用 Swagger 进行 API 文档管理
- Asynq 任务队列支持

## 安装

```bash
go mod tidy
```

## 配置

配置文件位于 `configs/` 目录中。项目支持配置中定义的多 MySQL 实例。

## 代码检查

本项目使用 golangci-lint 进行代码质量检查。

### 手动运行代码检查

```bash
# 使用提供的脚本
./scripts/lint.sh

# 或者直接运行
golangci-lint run ./...
```

### Git pre-commit 钩子

项目包含一个在每次提交前自动运行代码检查的 pre-commit 钩子。安装钩子：

```bash
./scripts/install-hooks.sh
```

钩子会在每次提交前自动运行，如果代码检查失败将阻止提交。要绕过检查（不推荐），使用：

```bash
git commit --no-verify
```

## 传输层

项目包含基于 Gin 框架的 HTTP 传输层：

- [HTTP 传输层文档](pkg/transport/rest/README.md)

## 消息队列

项目包含对 Asynq 的支持，Asynq 是一个简单、可靠且高效的 Go 分布式任务队列。Asynq 使用 Redis 作为消息代理，支持 Redis 的单机、集群和哨兵模式。

- [Asynq 客户端文档](pkg/client/mq/README.md)

## API 文档

本项目使用 Swagger 进行 API 文档管理。文档位于 `docs/swagger/swagger.json`。

### 初始化 Swagger 文档

```bash
make swagger-init
```

此命令创建 swagger 文档目录并确保其存在。

### 生成 Swagger 文档

```bash
make swagger-generate
```

此命令从源代码注释自动生成 Swagger 文档。文档从 `examples/` 目录中的示例文件生成。

### 验证 Swagger 文档

```bash
make swagger-validate
```

### 启动 Swagger UI

```bash
make swagger-serve
```

运行此命令后，打开浏览器并导航到 `http://localhost:63150/docs` 查看交互式 API 文档。

## 使用方法

各种组件的使用示例可以在 `examples/` 目录中找到。

## 文档

- [MySQL 多实例使用](pkg/client/database/README.md)
- [Redis 多实例使用](pkg/client/database/README.md)
- [Asynq 任务队列使用](pkg/client/mq/README.md)
- [HTTP 传输层使用](pkg/transport/rest/README.md)
- [配置管理](pkg/conf/README.md)
- [代码检查指南](docs/linting.md)