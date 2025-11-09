# Repository Layer Design

The repository layer in this service is designed to be replaceable, allowing different data storage implementations to be used without changing the business logic.

## Architecture

The repository layer follows a classic interface-implementation pattern:

1. **Interface**: Defined in `internal/repository/order.go`
2. **Implementations**:
   - In-memory storage: `internal/repository/order.go`
   - MySQL storage: `internal/repository/mysql/order.go`
   - MongoDB storage: `internal/repository/mongo/order.go`

## Interface Design

The `OrderRepository` interface defines all the operations that can be performed on orders:

```go
type OrderRepository interface {
    // Save saves an order.
    Save(order *model.Order) error

    // FindByID finds an order by ID.
    FindByID(id string) (*model.Order, error)

    // FindByUserID finds orders by user ID.
    FindByUserID(userID string) ([]*model.Order, error)

    // Update updates an order.
    Update(order *model.Order) error

    // Delete deletes an order by ID.
    Delete(id string) error
}
```

## Implementation Details

### In-Memory Repository

The in-memory implementation is useful for:
- Development and testing
- Simple use cases with low data volume
- Prototyping

It uses a map to store orders and a read-write mutex to ensure thread safety.

### MySQL Repository

The MySQL implementation provides:
- Persistent storage
- ACID compliance
- Scalability

Key features:
- Uses prepared statements to prevent SQL injection
- Handles JSON serialization for the items array
- Implements upsert functionality for Save operation
- Proper error handling and wrapping

### MongoDB Repository

The MongoDB implementation offers:
- Document-based storage
- Flexible schema
- Horizontal scaling

Key features:
- Uses MongoDB's native upsert functionality
- Leverages BSON for data serialization
- Implements proper context handling with timeouts
- Handles MongoDB-specific error cases

## Switching Between Implementations

The repository implementation is selected in `internal/server/server.go` based on the configuration:

```go
var repo repository.OrderRepository
switch cfg.Database.Driver {
case "mysql":
    // Initialize MySQL repository
    // repo = mysql.NewMySQLRepository(db)
    repo = repository.NewInMemoryOrderRepository()
case "mongo":
    // Initialize MongoDB repository
    // repo = mongo.NewMongoRepository(collection)
    repo = repository.NewInMemoryOrderRepository()
default:
    // Default to in-memory repository
    repo = repository.NewInMemoryOrderRepository()
}
```

Currently, the code is commented out for actual MySQL and MongoDB implementations, and uses the in-memory repository for demonstration purposes.

## Extending to Other Databases

To add support for other databases:

1. Create a new package under `internal/repository/` (e.g., `postgresql`)
2. Implement the `OrderRepository` interface
3. Update the switch statement in `internal/server/server.go` to handle the new driver type

## Best Practices

1. **Interface Segregation**: The repository interface is focused only on order operations
2. **Dependency Inversion**: Business logic depends on the interface, not concrete implementations
3. **Error Handling**: All implementations wrap errors with context for better debugging
4. **Thread Safety**: Implementations ensure thread safety where needed
5. **Resource Management**: Properly manage database connections and cursors