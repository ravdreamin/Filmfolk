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
	"filmfolk/internal/middleware"
	"filmfolk/internal/routes"
	"filmfolk/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration Loading Error: %v", err)
	}
	log.Printf("Configuration loaded from environment variables")
	log.Printf("Environment: %s", cfg.App.Env)

	// 2. Initialize structured logging
	log.Println("Initializing logger...")
	utils.InitLogger(cfg.App.Env)
	logger := utils.GetLogger()
	logger.Info().Msg("Logger initialized successfully")

	// 3. Initialize JWT utilities
	logger.Info().Msg("Initializing JWT...")
	utils.InitJWT(cfg.Jwt.Secret)

	// 4. Connect to database
	logger.Info().Msg("Connecting to database...")
	if err := db.InitDB(cfg); err != nil {
		logger.Fatal().Err(err).Msg("Database connection failed")
	}
	defer db.CloseDB()

	// 5. Auto-migrations disabled. Use the new migrate tool.
	// if cfg.App.Env == "development" {
	// 	logger.Info().Msg("Running auto-migrations...")
	// 	if err := db.AutoMigrate(); err != nil {
	// 		logger.Fatal().Err(err).Msg("Auto-migration failed")
	// 	}
	// }

	// 6. Set Gin mode based on environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 7. Create Gin router (without default middleware)
	router := gin.New()

	// 8. Add recovery middleware (panic recovery)
	router.Use(gin.Recovery())

	// 9. Add request ID middleware (must be first)
	router.Use(middleware.RequestIDMiddleware())

	// 10. Add structured logging middleware
	router.Use(middleware.LoggingMiddleware())

	// 11. Add security headers
	router.Use(middleware.SecurityHeadersMiddleware())

	// 12. Add CORS middleware with proper configuration
	router.Use(middleware.CORSMiddleware(cfg.App.AllowedOrigins, cfg.App.Env))

	// 13. Add global rate limiting
	router.Use(middleware.RateLimitMiddleware())

	// 14. Setup routes
	logger.Info().Msg("Setting up routes...")
	routes.SetupRoutes(router, cfg)

	// 15. Create HTTP server with production-ready timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 16. Start server in a goroutine
	go func() {
		logger.Info().
			Str("name", cfg.App.Name).
			Int("port", cfg.App.Port).
			Str("env", cfg.App.Env).
			Msgf("Server starting on http://localhost:%d", cfg.App.Port)

		logger.Info().Msgf("API endpoints: http://localhost:%d/api/v1", cfg.App.Port)
		logger.Info().Msgf("Health check: http://localhost:%d/health", cfg.App.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// 17. Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) or SIGTERM (Docker stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server gracefully...")

	// 18. Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited gracefully")
}
