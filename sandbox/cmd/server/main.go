package main

import (
	"log"

	"sandbox/internal/config"
	"sandbox/internal/docker"
	"sandbox/internal/handler"
	"sandbox/internal/router"
	"sandbox/internal/service"
)

// @title Sandbox API
// @version 1.0
// @description Extensible API server for managing Docker sandbox containers.
// @host localhost:8080
// @BasePath /
func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Initialize dependencies (Dependency Injection)
	dockerClient := docker.NewMockClient()
	sandboxService := service.NewSandboxService(dockerClient)
	sandboxHandler := handler.NewSandboxHandler(sandboxService)

	// 3. Initialize Router
	r := router.NewRouter(sandboxHandler)

	// 4. Start Server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
