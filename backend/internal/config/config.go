package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name           string   `mapstructure:"name" validate:"required"`
		Port           int      `mapstructure:"port" validate:"required,min=1,max=65535"`
		Env            string   `mapstructure:"env" validate:"required"` // development, production
		AllowedOrigins []string `mapstructure:"allowed_origins"`         // CORS allowed origins
	} `mapstructure:"app"`
	Db struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname" validate:"required"`
		SSLMode  string `mapstructure:"sslmode" validate:"required"`
	} `mapstructure:"db"`
	Jwt struct {
		Secret          string `mapstructure:"secret" validate:"required,min=16"`
		AccessTokenTTL  int    `mapstructure:"access_token_ttl" validate:"required"`  // minutes
		RefreshTokenTTL int    `mapstructure:"refresh_token_ttl" validate:"required"` // days
	} `mapstructure:"jwt"`
	OAuth struct {
		GoogleClientID     string `mapstructure:"google_client_id"`
		GoogleClientSecret string `mapstructure:"google_client_secret"`
		GoogleRedirectURL  string `mapstructure:"google_redirect_url"`
	} `mapstructure:"oauth"`
	TMDB struct {
		APIKey string `mapstructure:"api_key"` // For movie data
	} `mapstructure:"tmdb"`
	AI struct {
		OpenAIKey string `mapstructure:"openai_key"` // For content moderation & sentiment
	} `mapstructure:"ai"`
}

func LoadConfig() (*Config, error) {

	envConfig, err := loadEnvConfig()
	if err != nil {
		return nil, fmt.Errorf("environment configuration loading failed: %w", err)
	}

	if err := validateConfig(envConfig); err != nil {
		return nil, fmt.Errorf("environment configuration validation failed: %w", err)
	}

	return envConfig, nil
}

func loadEnvConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	v := viper.New()

	v.BindEnv("app.name", "APP_NAME")
	v.BindEnv("app.port", "APP_PORT")
	v.BindEnv("app.env", "APP_ENV")
	v.BindEnv("app.allowed_origins", "ALLOWED_ORIGINS")
	v.BindEnv("db.host", "DB_HOST")
	v.BindEnv("db.port", "DB_PORT")
	v.BindEnv("db.user", "DB_USER")
	v.BindEnv("db.password", "DB_PASSWORD")
	v.BindEnv("db.dbname", "DB_NAME")
	v.BindEnv("db.sslmode", "DB_SSLMODE")
	v.BindEnv("jwt.secret", "JWT_SECRET_KEY")
	v.BindEnv("jwt.access_token_ttl", "JWT_ACCESS_TOKEN_TTL")
	v.BindEnv("jwt.refresh_token_ttl", "JWT_REFRESH_TOKEN_TTL")
	v.BindEnv("oauth.google_client_id", "GOOGLE_CLIENT_ID")
	v.BindEnv("oauth.google_client_secret", "GOOGLE_CLIENT_SECRET")
	v.BindEnv("oauth.google_redirect_url", "GOOGLE_REDIRECT_URL")
	v.BindEnv("oauth.facebook_client_id", "FACEBOOK_CLIENT_ID")
	v.BindEnv("oauth.facebook_client_secret", "FACEBOOK_CLIENT_SECRET")
	v.BindEnv("oauth.facebook_redirect_url", "FACEBOOK_REDIRECT_URL")
	v.BindEnv("tmdb.api_key", "TMDB_API_KEY")
	v.BindEnv("ai.openai_key", "OPENAI_API_KEY")

	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling env config: %w", err)
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	var missingFields []string

	if cfg.App.Name == "" {
		missingFields = append(missingFields, "app.name")
	}
	if cfg.App.Port <= 0 || cfg.App.Port > 65535 {
		missingFields = append(missingFields, "app.port (must be between 1-65535)")
	}
	if cfg.App.Env == "" {
		missingFields = append(missingFields, "app.env")
	}

	if cfg.Db.Host == "" {
		missingFields = append(missingFields, "db.host")
	}
	if cfg.Db.Port <= 0 || cfg.Db.Port > 65535 {
		missingFields = append(missingFields, "db.port (must be between 1-65535)")
	}
	if cfg.Db.User == "" {
		missingFields = append(missingFields, "db.user")
	}
	if cfg.Db.DBName == "" {
		missingFields = append(missingFields, "db.dbname")
	}
	if cfg.Db.SSLMode == "" {
		missingFields = append(missingFields, "db.sslmode")
	}

	if cfg.Jwt.Secret == "" {
		missingFields = append(missingFields, "jwt.secret")
	} else if len(cfg.Jwt.Secret) < 16 {
		missingFields = append(missingFields, "jwt.secret (must be at least 16 characters)")
	}
	if cfg.Jwt.AccessTokenTTL <= 0 {
		missingFields = append(missingFields, "jwt.access_token_ttl (must be greater than 0)")
	}
	if cfg.Jwt.RefreshTokenTTL <= 0 {
		missingFields = append(missingFields, "jwt.refresh_token_ttl (must be greater than 0)")
	}

	if len(missingFields) > 0 {
		return errors.New("missing or invalid configuration fields: " + strings.Join(missingFields, ", "))
	}

	return nil
}

func GetConfigFilePath() string {
	configPath := "./configs/config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}
	return ""
}
