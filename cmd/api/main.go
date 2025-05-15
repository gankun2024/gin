package main

import (
	"cmp"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/gankun2024/gin-demo-project/internal/api/routes"
	"github.com/gankun2024/gin-demo-project/internal/config"
	"github.com/gankun2024/gin-demo-project/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize router
	r := gin.New()

	// Setup routes
	routes.Setup(r, log)

	// Get port from config or environment
	port := cmp.Or(os.Getenv("PORT"), cfg.Server.Port)

	log.Info("Server starting", slog.String("port", port))

	// Run server
	if err := r.Run(":" + port); err != nil {
		log.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
