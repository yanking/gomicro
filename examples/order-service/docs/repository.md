# Repository 层设计

本服务中的 repository 层被设计为可替换的，允许使用不同的数据存储实现而无需更改业务逻辑。

## 架构

repository 层遵循经典的接口-实现模式：

1. **接口**: 定义在 `internal/repository/order.go`
2. **实现**:
   - 内存存储: `internal/repository/order.go`
   - MySQL 存储: `internal/repository/mysql/order.go`
   - MongoDB 存储: `internal/repository/mongo/order.go`

## 接口设计

`OrderRepository` 接口定义了可以在订单上执行的所有操作：

```go
type OrderRepository interface {
    // Save 保存订单
    Save(order *model.Order) error

    // FindByID 根据 ID 查找订单
    FindByID(id string) (*model.Order, error)

    // FindByUserID 根据用户 ID 查找订单
    FindByUserID(userID string) ([]*model.Order, error)

    // Update 更新订单
    Update(order *model.Order) error

    // Delete 根据 ID 删除订单
    Delete(id string) error
}
```

## 实现详情

### 内存 Repository

内存实现实用于：
- 开发和测试
- 数据量较小的简单用例
- 原型开发

它使用 map 来存储订单，并使用读写互斥锁来确保线程安全。

### MySQL Repository

MySQL 实现提供：
- 持久化存储
- ACID 合规性
- 可扩展性

主要特性：
- 使用预处理语句防止 SQL 注入
- 处理项目数组的 JSON 序列化
- 为 Save 操作实现 upsert 功能
- 适当的错误处理和包装

### MongoDB Repository

MongoDB 实现提供：
- 基于文档的存储
- 灵活的模式
- 水平扩展

主要特性：
- 使用 MongoDB 的原生 upsert 功能
- 利用 BSON 进行数据序列化
- 实现带超时的适当上下文处理
- 处理 MongoDB 特定的错误情况

## 在实现之间切换

repository 实现基于配置在 `internal/server/server.go` 中选择：

```go
var repo repository.OrderRepository
switch cfg.Database.Driver {
case "mysql":
    // 初始化 MySQL repository
    // repo = mysql.NewMySQLRepository(db)
    repo = repository.NewInMemoryOrderRepository()
case "mongo":
    // 初始化 MongoDB repository
    // repo = mongo.NewMongoRepository(collection)
    repo = repository.NewInMemoryOrderRepository()
default:
    // 默认使用内存 repository
    repo = repository.NewInMemoryOrderRepository()
}
```

目前，实际的 MySQL 和 MongoDB 实现代码被注释掉了，使用内存 repository 作为演示用途。

## 扩展到其他数据库

要添加对其他数据库的支持：

1. 在 `internal/repository/` 下创建新包（例如，`postgresql`）
2. 实现 `OrderRepository` 接口
3. 更新 `internal/server/server.go` 中的 switch 语句以处理新的驱动类型

## 最佳实践

1. **接口隔离**: repository 接口仅专注于订单操作
2. **依赖倒置**: 业务逻辑依赖于接口，而不是具体实现
3. **错误处理**: 所有实现都用上下文包装错误以便更好地调试
4. **线程安全**: 实现在需要时确保线程安全
5. **资源管理**: 正确管理数据库连接和游标