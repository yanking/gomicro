# Logger Package

The logger package provides a simple wrapper around Go's `slog` package for consistent logging across the application.

## Features

- Easy initialization and configuration
- Support for different log formats (text, JSON)
- Support for different log levels 
- Global logger access
- Thread-safe operations
- Uses RFC3339 time format by default
- Optional source file and line number information
- Optional base path trimming for source file paths
- Automatic base path detection

## Usage

### Basic Usage

```go
import "github.com/yanking/gomicro/pkg/logger"

// Initialize the logger with default configuration
logger.Init(nil)

// Use the logger directly with slog methods
log := logger.Get()
log.Info("Application started")
log.Error("An error occurred", "error", err)
```

### Custom Configuration

```go
import "github.com/yanking/gomicro/pkg/logger"

config := &logger.Config{
    Level:              slog.LevelDebug,
    Format:             "json",
    Output:             os.Stdout,
    AddSource:          true,
    BasePath:           "/path/to/project/root/",
    AutoDetectBasePath: true,
}

logger.Init(config)
log := logger.Get()
log.Debug("Debug message")
```

### Direct Logger Instance

```go
import "github.com/yanking/gomicro/pkg/logger"

config := &logger.Config{
    Level:  slog.LevelInfo,
    Format: "text",
}

log := logger.New(config)
log.Info("Hello, world!")
```

## API

### Functions

- `Init(config *Config)` - Initialize the global logger
- `Get() *slog.Logger` - Get the global logger instance
- `DefaultConfig() *Config` - Get default configuration

### Types

- `Config` - Logger configuration

### Config Fields

- `Level` - Log level (slog.Level)
- `Format` - Log format ("text" or "json")
- `Output` - Output writer (default: os.Stdout)
- `AddSource` - Whether to add source file and line number
- `BasePath` - Base path to trim from source file paths
- `AutoDetectBasePath` - Automatically detect base path by looking for go.mod file