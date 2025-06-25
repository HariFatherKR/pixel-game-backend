package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		// Continue with defaults
	}

	// Initialize router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "pixel-game-backend",
			"timestamp": time.Now().Unix(),
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Placeholder endpoints
		api.GET("/cards", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"cards": []gin.H{
					{"id": 1, "name": "Code Slash", "type": "action", "cost": 2},
					{"id": 2, "name": "Firewall Up", "type": "action", "cost": 1},
					{"id": 3, "name": "Bug Found", "type": "event", "cost": 0},
				},
			})
		})

		api.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": "0.1.0",
				"build":   "dev",
			})
		})
	}

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}