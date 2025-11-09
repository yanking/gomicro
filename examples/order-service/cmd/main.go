package main

import (
	"flag"
	"os"

	appctx "github.com/yanking/gomicro/examples/order-service/internal/app"
	"github.com/yanking/gomicro/examples/order-service/internal/config"
	"github.com/yanking/gomicro/examples/order-service/internal/server"
	gomicroapp "github.com/yanking/gomicro/pkg/app"
	"github.com/yanking/gomicro/pkg/logger"
)

var configFile = flag.String("config", "config.yaml", "path to config file")

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		// Use default logger for config loading errors
		log := logger.Get()
		if log == nil {
			// If logger is not initialized, create a temporary one
			log = logger.New(logger.DefaultConfig())
		}
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize logger with config
	loggerConfig := &logger.Config{
		Level:              cfg.Log.GetLogLevel(),
		Format:             cfg.Log.Format,
		Output:             os.Stdout,
		AddSource:          cfg.Log.AddSource,
		BasePath:           cfg.Log.BasePath,
		AutoDetectBasePath: cfg.Log.AutoDetectBasePath,
	}
	logger.Init(loggerConfig)

	// Get logger instance
	log := logger.Get()

	// Create service context
	serviceContext := appctx.NewServiceContext(log, cfg)

	// Create server component
	orderServer := server.New(cfg)

	// Create application with components
	application, err := gomicroapp.New(serviceContext, "OrderService", "v1.0.0", orderServer)
	if err != nil {
		log.Error("failed to create application", "error", err)
		os.Exit(1)
	}

	// Run application
	if err := application.Run(); err != nil {
		log.Error("failed to run application", "error", err)
		os.Exit(1)
	}
}
