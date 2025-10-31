package main

import (
	"fmt"
	"log"
	"os"

	"filmfolk/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
)

func main() {
	// 1. Load configuration
	cfg, _, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration Loading Error: %v", err)
	}

	// 2. Build database connection string
	dsn := fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=%s",
		cfg.Db.User,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.DBName,
		cfg.Db.SSLMode,
	)

	// 3. Get migration files path
	// We assume the migrations are in the `migrations` directory relative to the project root.
	migrationsPath := "file://migrations"

	// 4. Create a new migrate instance
	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// 5. Get command-line arguments
	args := os.Args[1:]
	if len(args) == 0 {
		log.Println("Usage: go run cmd/migrate/main.go <command>")
		log.Println("Commands: up, down, drop, version")
		return
	}

	// 6. Execute the command
	cmd := args[0]
	switch cmd {
	case "up":
		log.Println("Running migrations up...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		log.Println("Migrations applied successfully.")
	case "down":
		log.Println("Running migrations down...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations down: %v", err)
		}
		log.Println("Migrations rolled back successfully.")
	case "drop":
		log.Println("Dropping database...")
		if err := m.Drop(); err != nil {
			log.Fatalf("Failed to drop database: %v", err)
		}
		log.Println("Database dropped successfully.")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		log.Printf("Version: %d, Dirty: %v", version, dirty)
	default:
		log.Printf("Unknown command: %s", cmd)
	}
}
