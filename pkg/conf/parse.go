// Package conf provides configuration parsing functionality.
// It supports parsing YAML configuration files and automatically binds environment variables.
//
// # Usage
//
// Basic usage:
//
//	type Config struct {
//	  Server struct {
//	    Host string `mapstructure:"host"`
//	    Port int    `mapstructure:"port"`
//	  } `mapstructure:"server"`
//	  Database struct {
//	    DSN string `mapstructure:"dsn"`
//	  } `mapstructure:"database"`
//	}
//
//	var cfg Config
//	if err := conf.Parse("config.yaml", &cfg); err != nil {
//	  log.Fatal(err)
//	}
//
//	// Access parsed values
//	fmt.Println("Server:", cfg.Server.Host, cfg.Server.Port)
//	fmt.Println("Database:", cfg.Database.DSN)
//
// Environment Variables:
// This package automatically binds environment variables. By default, it uses "GO_KIT" as
// the environment variable prefix and replaces dots with underscores.
//
// For example, with the struct above:
//   - GO_KIT_SERVER_HOST will bind to cfg.Server.Host
//   - GO_KIT_DATABASE_DSN will bind to cfg.Database.DSN
//
// Configuration Reload:
// The package supports automatic configuration reloading when the config file changes.
// Pass reload functions as additional arguments to enable this feature:
//
//	func reloadCallback() {
//	  // Handle configuration changes
//	  fmt.Println("Configuration reloaded")
//	}
//
//	if err := conf.Parse("config.yaml", &cfg, reloadCallback); err != nil {
//	  log.Fatal(err)
//	}
//
// In this case, whenever config.yaml is modified, reloadCallback will be executed.
package conf

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	mu sync.RWMutex
)

// Parse parses the configuration file into the provided struct.
// It supports YAML format and automatically binds environment variables.
//
// Parameters:
//   - configFile: path to the configuration file
//   - obj: pointer to the struct that will receive the parsed configuration
//   - reloads: optional functions that will be called when the configuration file changes
//
// Returns:
//   - error: if parsing fails
//
// Environment Variable Binding:
// The function automatically binds environment variables with the "GO_KIT" prefix.
// Dots in configuration keys are replaced with underscores in environment variable names.
// For example: server.host becomes GO_KIT_SERVER_HOST
func Parse(configFile string, obj any, reloads ...func()) error {
	// 创建独立的viper实例，避免全局实例带来的冲突
	v := viper.New()
	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configs file %s: %w", configFile, err)
	}

	v.AutomaticEnv()
	v.SetEnvPrefix("GO_KIT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	mu.Lock()
	err := v.Unmarshal(obj)
	mu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to unmarshal configs: %w", err)
	}

	if len(reloads) > 0 {
		watchConfig(v, obj, reloads...)
	}

	return nil
}

// watchConfig watches for configuration file changes and triggers reload callbacks.
func watchConfig(v *viper.Viper, obj any, reloads ...func()) {
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		mu.Lock()
		err := v.Unmarshal(obj)
		mu.Unlock()

		if err != nil {
			_ = fmt.Errorf("conf.watchConfig: viper.Unmarshal error: %v", err)
		} else {
			// 将 defer/recover 移到循环外面，对所有 reload 函数提供统一的 panic 保护
			defer func() {
				if r := recover(); r != nil {
					_ = fmt.Errorf("conf.watchConfig: reload function panic: %v", r)
				}
			}()

			for _, reload := range reloads {
				reload()
			}
		}
	})
}
