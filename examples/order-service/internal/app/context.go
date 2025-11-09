// Package app provides application context implementation.
package app

import (
	"log/slog"

	"github.com/yanking/gomicro/examples/order-service/internal/config"
	"github.com/yanking/gomicro/pkg/app"
)

// ConfigProvider implements the configuration provider interface.
type ConfigProvider struct {
	config *config.Config
}

// NewConfigProvider creates a new configuration provider.
func NewConfigProvider(config *config.Config) *ConfigProvider {
	return &ConfigProvider{
		config: config,
	}
}

// GetConfig returns the configuration.
func (c *ConfigProvider) GetConfig() *config.Config {
	return c.config
}

// Ensure ConfigProvider implements app.IConfigProvider
var _ app.IConfigProvider = (*ConfigProvider)(nil)

// ServiceContext implements the service context interface.
type ServiceContext struct {
	logger *slog.Logger
	config *config.Config
}

// NewServiceContext creates a new service context.
func NewServiceContext(logger *slog.Logger, config *config.Config) *ServiceContext {
	return &ServiceContext{
		logger: logger,
		config: config,
	}
}

// GetLogger returns the logger.
func (s *ServiceContext) GetLogger() *slog.Logger {
	return s.logger
}

// GetConfig returns the configuration provider.
func (s *ServiceContext) GetConfig() app.IConfigProvider {
	return NewConfigProvider(s.config)
}
