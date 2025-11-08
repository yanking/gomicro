# Order Service

基于Go微服务最佳实践的订单服务示例。

## 项目结构

```
.
├── api/                    # API定义（如Swagger/OpenAPI）
├── cmd/                    # 应用程序入口点
│   └── main.go            # 主程序入口
├── configs/                # 配置文件
├── deploy/                 # 部署相关文件（Dockerfile, Kubernetes配置等）
├── docs/                   # 文档
├── internal/               # 私有应用代码
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP处理程序
│   ├── middleware/        # 中间件
│   ├── model/             # 领域模型
│   ├── repository/        # 数据访问层
│   ├── server/            # HTTP服务器封装
│   └── service/           # 业务逻辑层
├── scripts/                # 脚本文件
└── go.mod                 # Go模块定义
```

## 架构说明

本项目遵循Go微服务最佳实践，采用标准的分层架构：

1. **cmd/** - 应用程序入口，负责初始化和启动服务
2. **internal/** - 私有应用代码，外部无法直接引用
   - **config/** - 配置管理
   - **model/** - 领域模型（DDD中的实体和值对象）
   - **repository/** - 数据访问层，封装数据存储细节
   - **service/** - 业务逻辑层，实现核心业务规则
   - **handler/** - HTTP处理程序，处理HTTP请求和响应
   - **server/** - HTTP服务器封装，负责路由注册和服务器生命周期管理
   - **middleware/** - 中间件组件（如日志、认证、限流等）

## 运行服务

```bash
# 进入项目目录
cd examples/order-service

# 运行服务
go run cmd/main.go
```

服务将在 `localhost:8080` 启动。

## API接口

### 创建订单
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order-1",
    "user_id": "user-1",
    "items": [
      {
        "product_id": "product-1",
        "name": "iPhone 13",
        "price": 999.99,
        "quantity": 1
      }
    ]
  }'
```

### 获取订单
```bash
curl http://localhost:8080/orders/get?id=order-1
```

### 支付订单
```bash
curl -X POST http://localhost:8080/orders/pay?id=order-1
```

### 发货订单
```bash
curl -X POST http://localhost:8080/orders/ship?id=order-1
```

### 送达订单
```bash
curl -X POST http://localhost:8080/orders/deliver?id=order-1
```

### 取消订单
```bash
curl -X POST http://localhost:8080/orders/cancel?id=order-1
```

## 最佳实践说明

### 1. 项目结构
- 使用 `internal/` 目录存放私有代码，防止外部直接引用
- 按功能模块划分目录结构，便于维护和扩展
- 遵循Go社区推荐的项目布局

### 2. 分层架构
- **Handler层**：处理HTTP请求和响应，数据转换
- **Service层**：实现业务逻辑，处理领域规则
- **Repository层**：封装数据访问细节
- **Model层**：定义领域模型和实体

### 3. 依赖管理
- 通过接口定义抽象，实现依赖倒置
- 依赖通过构造函数注入，便于测试和替换

### 4. 配置管理
- 集中化配置管理，支持多种配置源（文件、环境变量等）
- 结构化配置定义，便于维护

### 5. 服务生命周期管理
- 优雅启动和关闭服务
- 支持上下文控制和超时处理

## 扩展建议

1. **数据库集成**：替换内存存储为真实的数据库（MySQL、PostgreSQL等）
2. **日志系统**：集成结构化日志系统（如zap、logrus）
3. **监控指标**：添加Prometheus指标收集
4. **链路追踪**：集成OpenTelemetry实现分布式追踪
5. **配置中心**：集成配置中心（如Consul、Etcd）
6. **服务注册发现**：集成服务注册发现机制
7. **限流熔断**：添加限流和熔断机制
8. **认证授权**：实现JWT认证和RBAC权限控制