// Package config provides configuration management for the service.
package config

// Config holds the application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig holds the server configuration.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

// Load loads the configuration from the environment or config file.
func Load() (*Config, error) {
	// For simplicity, we're using default values here.
	// In a real application, you would load from a file or environment variables.
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Database: DatabaseConfig{
			Driver:   "memory",
			Host:     "localhost",
			Port:     "3306",
			Username: "user",
			Password: "password",
			Name:     "orderdb",
		},
	}

	return cfg, nil
}
