package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/auth"
	"github.com/yourusername/pixel-game/internal/config"
	"github.com/yourusername/pixel-game/internal/database"
	"github.com/yourusername/pixel-game/internal/handlers"
	"github.com/yourusername/pixel-game/internal/middleware"
	"github.com/yourusername/pixel-game/internal/repository/postgres"
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

	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepository := postgres.NewUserRepository(db.DB)
	cardRepository := postgres.NewCardRepository(db.DB)

	// Initialize JWT manager
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		jwtSecretKey = "your-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret key. Set JWT_SECRET_KEY environment variable for production.")
	}
	jwtManager := auth.NewJWTManager(jwtSecretKey)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(jwtManager, userRepository, cardRepository)
	userHandler := handlers.NewUserHandler(userRepository)
	cardHandler := handlers.NewCardHandler(cardRepository, jwtManager)

	// Initialize router
	r := gin.Default()

	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",  // React default
			"http://localhost:5173",  // Vite default
			"http://localhost:8080",  // Same origin
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:8080",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Content-Length", "Accept-Encoding",
			"X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup Swagger
	swagger.SetupSwagger(r)

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check endpoint
		api.GET("/health", HealthCheck)
		
		// Authentication endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			
			// Protected auth endpoints
			auth.Use(middleware.AuthMiddleware(jwtManager))
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/profile", authHandler.Profile)
		}
		
		// User management endpoints (protected)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtManager))
		{
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/stats", userHandler.GetStats)
			users.GET("/collection", userHandler.GetCollection)
			users.POST("/stats/games-played", userHandler.IncrementGamesPlayed)
			users.POST("/stats/games-won", userHandler.IncrementGamesWon)
			users.POST("/stats/play-time/:seconds", userHandler.AddPlayTime)
		}
		
		// Card endpoints
		cardHandler.RegisterRoutes(api)
		
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