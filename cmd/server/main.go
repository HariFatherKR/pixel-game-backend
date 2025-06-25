package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hypertokki/pixel-game/internal/config"
	"github.com/hypertokki/pixel-game/internal/infrastructure/cache"
	"github.com/hypertokki/pixel-game/internal/infrastructure/persistence"
	"github.com/hypertokki/pixel-game/internal/interfaces/http/handler"
	"github.com/hypertokki/pixel-game/internal/interfaces/http/middleware"
	"github.com/hypertokki/pixel-game/internal/interfaces/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := persistence.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Set Gin mode
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg.CORS))

	// Initialize handlers
	authHandler := handler.NewAuthHandler(db, redisClient, cfg.JWT)
	userHandler := handler.NewUserHandler(db)
	cardHandler := handler.NewCardHandler(db)
	gameHandler := handler.NewGameHandler(db, redisClient)

	// Setup routes
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", middleware.JWTAuth(cfg.JWT), authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWT))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.GET("/stats", userHandler.GetStats)
				users.GET("/collection", userHandler.GetCollection)
			}

			// Card routes
			cards := protected.Group("/cards")
			{
				cards.GET("", cardHandler.ListCards)
				cards.GET("/:id", cardHandler.GetCard)
			}

			// Game routes
			games := protected.Group("/games")
			{
				games.POST("/start", gameHandler.StartGame)
				games.GET("/:id", gameHandler.GetGameState)
				games.POST("/:id/actions", gameHandler.PerformAction)
				games.POST("/:id/end", gameHandler.EndGame)
			}

			// Leaderboard
			protected.GET("/leaderboard", gameHandler.GetLeaderboard)

			// Daily challenge
			protected.GET("/challenges/daily", gameHandler.GetDailyChallenge)
		}
	}

	// WebSocket endpoint
	wsHub := websocket.NewHub(redisClient)
	go wsHub.Run()
	
	router.GET("/ws", middleware.JWTWebSocket(cfg.JWT), func(c *gin.Context) {
		websocket.ServeWS(wsHub, c.Writer, c.Request, c.GetString("userID"))
	})

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}