// Package config provides configuration management for the service.
package config

import (
	"fmt"
	"log/slog"

	"github.com/yanking/gomicro/pkg/conf"
)

// Config holds the application configuration.
type Config struct {
	Server   ServerConfig     `mapstructure:"server"`
	Log      LogConfig        `mapstructure:"log"`
	Database []DatabaseConfig `mapstructure:"database"`
}

// ServerConfig holds the server configuration.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// LogConfig holds the logger configuration.
type LogConfig struct {
	Level              string `mapstructure:"level"`                 // "debug", "info", "warn", "error"
	Format             string `mapstructure:"format"`                // "text", "json"
	AddSource          bool   `mapstructure:"add_source"`            // whether to add source file and line number
	BasePath           string `mapstructure:"base_path"`             // base path to trim from source file paths
	AutoDetectBasePath bool   `mapstructure:"auto_detect_base_path"` // automatically detect base path
}

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	Instance string `mapstructure:"instance"`
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

// Load loads the configuration from the environment or config file.
func Load(configFile string) (*Config, error) {
	var cfg Config

	// If configFile is provided, use it directly
	if configFile != "" {
		if err := conf.Parse(configFile, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file %s: %w", configFile, err)
		}
		return &cfg, nil
	}

	// Try to load from default config file location
	if err := conf.Parse("./configs/config.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse default config file: %w", err)
	}

	return &cfg, nil
}

// GetDatabaseConfig returns the database configuration for the specified instance.
func (c *Config) GetDatabaseConfig(instance string) *DatabaseConfig {
	for _, db := range c.Database {
		if db.Instance == instance {
			return &db
		}
	}
	return nil
}

// GetDefaultDatabaseConfig returns the default database configuration.
func (c *Config) GetDefaultDatabaseConfig() *DatabaseConfig {
	return c.GetDatabaseConfig("default")
}

// GetLogLevel returns the slog.Level based on the Level string.
func (l *LogConfig) GetLogLevel() slog.Level {
	switch l.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
