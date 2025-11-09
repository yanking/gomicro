// Package logger provides a simple wrapper around slog for consistent logging across the application.
package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	// defaultLogger is the default logger instance
	defaultLogger *slog.Logger
	// loggerMutex protects defaultLogger during initialization
	loggerMutex sync.RWMutex
	// basePath is used to trim the base path from source file paths
	basePath string
)

// Config holds logger configuration
type Config struct {
	// Level is the logging level
	Level slog.Level
	// Format is the log format: "text" or "json"
	Format string
	// Output is the log output writer, default to os.Stdout
	Output io.Writer
	// AddSource determines whether to add source file and line number
	AddSource bool
	// BasePath is the base path to trim from source file paths
	BasePath string
	// AutoDetectBasePath determines whether to automatically detect the base path
	AutoDetectBasePath bool
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:              slog.LevelInfo,
		Format:             "text",
		Output:             os.Stdout,
		AddSource:          false,
		BasePath:           "",
		AutoDetectBasePath: true,
	}
}

// detectBasePath tries to automatically detect the project base path
func detectBasePath() string {
	// Get the caller's file path
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	// Try to find a reasonable base path
	// Usually, we want to trim up to the project root
	dir := filepath.Dir(file)

	// Walk up the directory tree to find a marker of project root
	// Common markers: go.mod file
	for {
		// Check if go.mod exists in this directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found the project root
			return dir + string(filepath.Separator)
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the filesystem root, stop here
			break
		}
		dir = parent
	}

	// If we couldn't find go.mod, try another approach
	// Use the directory of the main module as base path
	return filepath.Dir(file)
}

// New creates a new logger instance based on the provided configuration
func New(config *Config) *slog.Logger {
	if config == nil {
		config = DefaultConfig()
	}

	if config.Output == nil {
		config.Output = os.Stdout
	}

	// Set the base path for trimming
	if config.BasePath != "" {
		basePath = config.BasePath
	} else if config.AutoDetectBasePath {
		basePath = detectBasePath()
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	// Custom attribute replacement function
	opts.ReplaceAttr = func(_ []string, a slog.Attr) slog.Attr {
		// Use RFC3339 time format by default
		if a.Key == slog.TimeKey {
			if t, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(t.Format(time.RFC3339))
			}
		}

		// Trim base path from source file paths
		if a.Key == slog.SourceKey && basePath != "" {
			if source, ok := a.Value.Any().(*slog.Source); ok && source != nil {
				if strings.HasPrefix(source.File, basePath) {
					source.File = filepath.Join("./", strings.TrimPrefix(source.File, basePath))
					// Remove leading slash if present
					source.File = strings.TrimPrefix(source.File, "/")
					// Clean up the path
					source.File = filepath.Clean(source.File)
				}
			}
		}

		return a
	}

	if strings.ToLower(config.Format) == "json" {
		handler = slog.NewJSONHandler(config.Output, opts)
	} else {
		handler = slog.NewTextHandler(config.Output, opts)
	}

	return slog.New(handler)
}

// Init initializes the default logger with the provided configuration
func Init(config *Config) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	defaultLogger = New(config)
}

// Get returns the default logger instance
// If the default logger hasn't been initialized, it creates one with default configuration
func Get() *slog.Logger {
	loggerMutex.RLock()
	logger := defaultLogger
	loggerMutex.RUnlock()

	if logger == nil {
		loggerMutex.Lock()
		if defaultLogger == nil {
			defaultLogger = New(DefaultConfig())
		}
		logger = defaultLogger
		loggerMutex.Unlock()
	}

	return logger
}
