package db

import (
	"fmt"
	"log"
	"time"

	"filmfolk/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection
// We make it a package-level variable so it's accessible everywhere
var DB *gorm.DB

// InitDB initializes the database connection and runs migrations
// This is called at application startup
func InitDB(cfg *config.Config) error {
	// Build PostgreSQL connection string (DSN - Data Source Name)
	// Format: "host=X user=Y password=Z dbname=W port=N sslmode=disable"
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%d sslmode=%s",
		cfg.Db.Host,
		cfg.Db.User,
		cfg.Db.DBName,
		cfg.Db.Port,
		cfg.Db.SSLMode,
	)

	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if cfg.App.Env == "development" {
		// In development: Show all SQL queries for debugging
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		// In production: Only log errors to reduce noise
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection with configuration
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			// Use UTC for all timestamps
			// Why? Avoids timezone confusion in distributed systems
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB to configure connection pooling
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	// Why? Reusing connections is MUCH faster than creating new ones
	sqlDB.SetMaxIdleConns(10)           // Keep 10 connections open even when idle
	sqlDB.SetMaxOpenConns(100)          // Max 100 concurrent connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Close connections after 1 hour

	log.Println("Database connection established successfully")
	return nil
}

// CloseDB closes the database connection
// Call this on application shutdown
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
