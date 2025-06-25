package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yourusername/pixel-game/docs"
)

// @title Pixel Game API
// @version 1.0
// @description API server for Cyberpunk Deck Building Card Game
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.pixelgame.io/support
// @contact.email support@pixelgame.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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
	r.GET("/health", HealthCheck)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	{
		// Card endpoints
		api.GET("/cards", GetCards)
		api.GET("/cards/:id", GetCard)
		
		// Version endpoint
		api.GET("/version", GetVersion)
	}

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Service:   "pixel-game-backend",
		Timestamp: time.Now().Unix(),
	})
}

// GetCards godoc
// @Summary List all cards
// @Description Get a list of all available cards
// @Tags cards
// @Accept json
// @Produce json
// @Success 200 {object} CardsResponse
// @Router /cards [get]
func GetCards(c *gin.Context) {
	c.JSON(http.StatusOK, CardsResponse{
		Cards: []Card{
			{ID: 1, Name: "Code Slash", Type: "action", Cost: 2, Description: "Deal 8 damage and apply Vulnerable"},
			{ID: 2, Name: "Firewall Up", Type: "action", Cost: 1, Description: "Gain 10 Shield"},
			{ID: 3, Name: "Bug Found", Type: "event", Cost: 0, Description: "Disable all traps, gain 1 random card"},
		},
		Total: 3,
	})
}

// GetCard godoc
// @Summary Get card by ID
// @Description Get details of a specific card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} Card
// @Failure 404 {object} ErrorResponse
// @Router /cards/{id} [get]
func GetCard(c *gin.Context) {
	// This is a placeholder implementation
	c.JSON(http.StatusOK, Card{
		ID:          1,
		Name:        "Code Slash",
		Type:        "action",
		Cost:        2,
		Description: "Deal 8 damage and apply Vulnerable",
	})
}

// GetVersion godoc
// @Summary Get API version
// @Description Get the current version of the API
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} VersionResponse
// @Router /version [get]
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, VersionResponse{
		Version: "0.1.0",
		Build:   "dev",
	})
}

// Response Models

type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Service   string `json:"service" example:"pixel-game-backend"`
	Timestamp int64  `json:"timestamp" example:"1234567890"`
}

type Card struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Code Slash"`
	Type        string `json:"type" example:"action" enums:"action,event,power"`
	Cost        int    `json:"cost" example:"2"`
	Description string `json:"description" example:"Deal 8 damage and apply Vulnerable"`
}

type CardsResponse struct {
	Cards []Card `json:"cards"`
	Total int    `json:"total" example:"3"`
}

type VersionResponse struct {
	Version string `json:"version" example:"0.1.0"`
	Build   string `json:"build" example:"dev"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Not found"`
	Message string `json:"message" example:"The requested resource was not found"`
}