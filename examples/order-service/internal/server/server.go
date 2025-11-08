// Package server provides HTTP server functionality.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yanking/gomicro/examples/order-service/internal/config"
	"github.com/yanking/gomicro/examples/order-service/internal/handler"
	"github.com/yanking/gomicro/examples/order-service/internal/repository"
	"github.com/yanking/gomicro/examples/order-service/internal/service"
)

// Server wraps the HTTP server and application components.
type Server struct {
	httpServer *http.Server
	repo       repository.OrderRepository
	service    *service.OrderService
	handler    *handler.OrderHandler
	config     *config.Config
}

// New creates a new Server instance.
func New(cfg *config.Config) *Server {
	// Initialize components
	repo := repository.NewInMemoryOrderRepository()
	orderService := service.NewOrderService(repo)
	orderHandler := handler.NewOrderHandler(orderService)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr: addr,
	}

	return &Server{
		httpServer: httpServer,
		repo:       repo,
		service:    orderService,
		handler:    orderHandler,
		config:     cfg,
	}
}

// RegisterRoutes registers the HTTP routes.
func (s *Server) RegisterRoutes() {
	http.HandleFunc("/orders", s.handler.CreateOrder)
	http.HandleFunc("/orders/get", s.handler.GetOrder)
	http.HandleFunc("/orders/pay", s.handler.PayOrder)
	http.HandleFunc("/orders/ship", s.handler.ShipOrder)
	http.HandleFunc("/orders/deliver", s.handler.DeliverOrder)
	http.HandleFunc("/orders/cancel", s.handler.CancelOrder)
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	s.RegisterRoutes()

	log.Printf("Starting server on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// WaitForShutdown waits for shutdown signal and gracefully stops the server.
func (s *Server) WaitForShutdown() {
	// In a real application, you would listen for OS signals like SIGINT or SIGTERM
	// For simplicity, we'll just wait for a short duration in this example
	time.Sleep(10 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		log.Printf("Error shutting down server: %v\n", err)
	}
}
