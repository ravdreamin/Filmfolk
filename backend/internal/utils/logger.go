package utils

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger initializes the global structured logger
func InitLogger(env string) {
	var output io.Writer = os.Stdout

	// Configure based on environment
	if env == "development" {
		// Pretty, colorized output for development
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		// JSON output for production (for log aggregation tools)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Create logger with timestamp
	Logger = zerolog.New(output).With().
		Timestamp().
		Caller().
		Str("service", "filmfolk").
		Logger()

	// Set as global logger
	log.Logger = Logger
}

// GetLogger returns the global logger
func GetLogger() *zerolog.Logger {
	return &Logger
}
