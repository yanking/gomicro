# Configuration Management

The order service uses a strict configuration loading strategy. If the configuration file cannot be loaded or parsed, the service will fail to start. This ensures that the service always runs with the correct configuration.

## Command-Line Flags

The service accepts a command-line flag to specify the configuration file path:

```bash
-config string
    Path to the configuration file (default "./configs/config.yaml")
```

Example usage:
```bash
# Use default config file location
./order-service

# Specify a custom config file location
./order-service -config /path/to/custom/config.yaml
```

## Configuration File Structure

The configuration file follows this structure:

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

### Server Configuration

The `server` section defines the HTTP server settings:
- `port`: The port on which the server listens
- `host`: The host address on which the server binds

### Database Configuration

The `database` section is a list of database configurations, supporting multiple instances:
- `instance`: The name of the database instance (e.g., "default", "mysql", "mongo")
- `driver`: The database driver to use ("memory", "mysql", "mongo")
- `host`: The database host (for MySQL)
- `port`: The database port (for MySQL)
- `username`: The database username (for MySQL)
- `password`: The database password (for MySQL)
- `name`: The database name (for MySQL)
- `uri`: The connection URI (for MongoDB)
- `database`: The database name (for MongoDB)

## Using Configuration in Code

The configuration is loaded at startup using the `config.Load()` function:

```go
cfg, err := config.Load(*configFile)
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}
```

To get the default database configuration:

```go
defaultDBConfig := cfg.GetDefaultDatabaseConfig()
```

To get a specific database instance configuration:

```go
mysqlConfig := cfg.GetDatabaseConfig("mysql")
```

## Environment Variables

The service also supports environment variables through the `conf` package. Any configuration value can be overridden by setting an environment variable with the prefix `GO_KIT_` and with dots replaced by underscores.

For example:
- `GO_KIT_SERVER_PORT` overrides `server.port`
- `GO_KIT_DATABASE_0_USERNAME` overrides the username of the first database instance

## Error Handling

If the configuration file cannot be loaded or parsed, the service will exit with an error message. This strict behavior ensures that the service always runs with a valid configuration.

## Best Practices

1. **Always Provide Configuration**: Ensure a valid configuration file is always available when starting the service.

2. **Validate Configuration**: Check that all required configuration values are present and valid.

3. **Environment-Specific Configuration**: Use different configuration files for different environments (development, staging, production).

4. **Sensitive Information**: Never commit sensitive information like passwords to version control. Use environment variables or secure configuration management systems for sensitive data.

5. **Documentation**: Keep the configuration file well-documented so that other developers understand the available options.