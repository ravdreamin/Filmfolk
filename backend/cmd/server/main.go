package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"filmfolk/internal/config"
	"filmfolk/internal/db"
	"filmfolk/internal/routes"
	"filmfolk/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	log.Println("Loading configuration...")
	cfg, source, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration Loading Error: %v", err)
	}
	log.Printf("Configuration loaded from: %s", source)
	log.Printf("Environment: %s", cfg.App.Env)

	// 2. Initialize JWT utilities
	log.Println("Initializing JWT...")
	utils.InitJWT(cfg.Jwt.Secret)

	// 3. Connect to database
	log.Println("Connecting to database...")
	if err := db.InitDB(cfg); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.CloseDB()

	// 4. Run auto-migrations (development only!)
	if cfg.App.Env == "development" {
		log.Println("Running auto-migrations...")
		if err := db.AutoMigrate(); err != nil {
			log.Fatalf("Auto-migration failed: %v", err)
		}
	}

	// 5. Set Gin mode based on environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 6. Create Gin router
	router := gin.Default()

	// 7. Setup CORS middleware (allow frontend to access API)
	router.Use(corsMiddleware())

	// 8. Setup routes
	log.Println("Setting up routes...")
	routes.SetupRoutes(router, cfg)

	// 9. Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 10. Start server in a goroutine
	go func() {
		log.Printf("🚀 %s server starting on http://localhost:%d", cfg.App.Name, cfg.App.Port)
		log.Printf("📝 API docs will be at http://localhost:%d/api/v1", cfg.App.Port)
		log.Printf("❤️  Health check: http://localhost:%d/health", cfg.App.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// 11. Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) or SIGTERM (Docker stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// 12. Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited gracefully")
}

// corsMiddleware handles Cross-Origin Resource Sharing
// Allows frontend apps on different domains to access the API
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
