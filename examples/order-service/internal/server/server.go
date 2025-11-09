// Package server provides HTTP server functionality using the project's REST framework.
package server

import (
	"context"
	"log/slog"

	"github.com/yanking/gomicro/examples/order-service/internal/config"
	"github.com/yanking/gomicro/examples/order-service/internal/handler"
	"github.com/yanking/gomicro/examples/order-service/internal/repository"
	"github.com/yanking/gomicro/examples/order-service/internal/service"
	"github.com/yanking/gomicro/pkg/lifecycle"
	"github.com/yanking/gomicro/pkg/logger"
	"github.com/yanking/gomicro/pkg/transport/rest"
)

// Server wraps the REST server and application components.
type Server struct {
	restServer *rest.Server
	repo       repository.OrderRepository
	service    *service.OrderService
	handler    *handler.OrderHandler
	config     *config.Config
	logger     *slog.Logger
}

// New creates a new Server instance.
func New(cfg *config.Config) *Server {
	// Get logger instance
	log := logger.Get()

	// Initialize components
	// For now, we're using the in-memory repository
	// In a production environment, you would initialize
	// the appropriate repository based on configuration
	var repo repository.OrderRepository

	// Get the default database configuration
	defaultDBConfig := cfg.GetDefaultDatabaseConfig()
	if defaultDBConfig != nil {
		switch defaultDBConfig.Driver {
		case "mysql":
			// Initialize MySQL repository
			// repo = mysql.NewMySQLRepository(db)
			// For demonstration purposes, we're still using in-memory
			repo = repository.NewInMemoryOrderRepository()
			log.Info("Using MySQL repository")
		case "mongo":
			// Initialize MongoDB repository
			// repo = mongo.NewMongoRepository(collection)
			// For demonstration purposes, we're still using in-memory
			repo = repository.NewInMemoryOrderRepository()
			log.Info("Using MongoDB repository")
		default:
			// Default to in-memory repository
			repo = repository.NewInMemoryOrderRepository()
			log.Info("Using in-memory repository (default)")
		}
	} else {
		// Fallback to in-memory repository if no database config is found
		repo = repository.NewInMemoryOrderRepository()
		log.Info("Using in-memory repository (fallback)")
	}

	orderService := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(orderService)

	// Create REST server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	restServer := rest.NewServer(log, rest.WithAddr(addr))

	// Register routes
	registerRoutes(restServer, orderHandler)

	return &Server{
		restServer: restServer,
		repo:       repo,
		service:    orderService,
		handler:    orderHandler,
		config:     cfg,
		logger:     log,
	}
}

// registerRoutes registers the HTTP routes.
func registerRoutes(server *rest.Server, handler *handler.OrderHandler) {
	server.POST("/orders", handler.CreateOrder)
	server.GET("/orders/get", handler.GetOrder)
	server.POST("/orders/pay", handler.PayOrder)
	server.POST("/orders/ship", handler.ShipOrder)
	server.POST("/orders/deliver", handler.DeliverOrder)
	server.POST("/orders/cancel", handler.CancelOrder)
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting order service server", "addr", s.config.Server.Host+":"+s.config.Server.Port)
	return s.restServer.Start(ctx)
}

// Stop gracefully stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping order service server")
	return s.restServer.Stop(ctx)
}

// Name returns the component name.
func (s *Server) Name() string {
	return "order-service-server"
}

// Ensure Server implements lifecycle.Component interface
var _ lifecycle.Component = (*Server)(nil)
