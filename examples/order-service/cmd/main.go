// Package main is the entry point for the order service.
package main

import (
	"log"

	"github.com/yanking/gomicro/examples/order-service/internal/config"
	"github.com/yanking/gomicro/examples/order-service/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	s := server.New(cfg)

	// Start server
	log.Println("Starting order service...")
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
