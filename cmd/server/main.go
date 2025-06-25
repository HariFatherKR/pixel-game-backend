package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/config"
	"github.com/yourusername/pixel-game/internal/swagger"
	_ "github.com/yourusername/pixel-game/docs"
)

// @title           Pixel Game - 사이버펑크 덱 빌딩 카드 게임 API
// @version         0.1.0
// @description     Vibe 코딩 기반 사이버펑크 덱 빌딩 카드 게임의 백엔드 API 서버입니다. 프론트엔드 코드 로직과 게임플레이를 연결하여 카드 사용 시 실제 코드가 실행되는 혁신적인 게임입니다.
// @termsOfService  https://github.com/HariFatherKR/pixel-game-backend

// @contact.name   Pixel Game Backend Team
// @contact.url    https://github.com/HariFatherKR/pixel-game-backend/issues
// @contact.email  support@pixelgame.io

// @license.name  MIT License
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		// Continue with defaults
	}

	// Initialize router
	r := gin.Default()

	// Setup Swagger
	swagger.SetupSwagger(r)

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check endpoint
		api.GET("/health", HealthCheck)
		
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
// @Summary      Health check
// @Description  Check if the service is healthy
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Service:   "pixel-game-backend",
		Timestamp: time.Now().Unix(),
	})
}

// GetCards godoc
// @Summary      카드 목록 조회
// @Description  게임에서 사용 가능한 모든 카드 목록을 조회합니다. 각 카드는 고유한 코드 효과를 가지고 있습니다.
// @Tags         cards
// @Accept       json
// @Produce      json
// @Success      200  {object}  CardsResponse  "카드 목록"
// @Router       /cards [get]
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
// @Summary      카드 상세 조회
// @Description  특정 카드의 상세 정보를 조회합니다. 카드의 코드 효과와 시각적 효과 정보를 포함합니다.
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "카드 ID"
// @Success      200  {object}  Card          "카드 상세 정보"
// @Failure      404  {object}  ErrorResponse "카드를 찾을 수 없음"
// @Router       /cards/{id} [get]
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
// @Summary      Get API version
// @Description  Get the current version of the API
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  VersionResponse
// @Router       /version [get]
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, VersionResponse{
		Version: "0.1.0",
		Build:   "dev",
	})
}

// Response Models

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Service   string `json:"service" example:"pixel-game-backend"`
	Timestamp int64  `json:"timestamp" example:"1234567890"`
}

// Card represents a game card
type Card struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Code Slash"`
	Type        string `json:"type" example:"action" enums:"action,event,power"`
	Cost        int    `json:"cost" example:"2"`
	Description string `json:"description" example:"Deal 8 damage and apply Vulnerable"`
}

// CardsResponse represents the cards list response
type CardsResponse struct {
	Cards []Card `json:"cards"`
	Total int    `json:"total" example:"3"`
}

// VersionResponse represents the version info response
type VersionResponse struct {
	Version string `json:"version" example:"0.1.0"`
	Build   string `json:"build" example:"dev"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Not found"`
	Message string `json:"message" example:"The requested resource was not found"`
}