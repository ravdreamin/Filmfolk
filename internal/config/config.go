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
		Name string `mapstructure:"name" validate:"required"`
		Port int    `mapstructure:"port" validate:"required,min=1,max=65535"`
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
		Secret string `mapstructure:"secret" validate:"required,min=16"`
	} `mapstructure:"jwt"`
}

type ConfigSource string

const (
	SourceYAML ConfigSource = "YAML"
	SourceEnv  ConfigSource = "Environment Variables"
)

func LoadConfig() (*Config, ConfigSource, error) {
	yamlConfig, err := loadYamlConfig()
	if err == nil && yamlConfig != nil {
		if err := validateConfig(yamlConfig); err != nil {
			return nil, SourceYAML, fmt.Errorf("YAML configuration validation failed: %w", err)
		}
		return yamlConfig, SourceYAML, nil
	}

	envConfig, err := loadEnvConfig()
	if err != nil {
		return nil, SourceEnv, fmt.Errorf("environment configuration loading failed: %w", err)
	}

	if err := validateConfig(envConfig); err != nil {
		return nil, SourceEnv, fmt.Errorf("environment configuration validation failed: %w", err)
	}

	return envConfig, SourceEnv, nil
}

func loadYamlConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read YAML config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML config: %w", err)
	}

	return &cfg, nil
}

func loadEnvConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	v := viper.New()

	v.BindEnv("app.name", "APP_NAME")
	v.BindEnv("app.port", "APP_PORT")
	v.BindEnv("db.host", "DB_HOST")
	v.BindEnv("db.port", "DB_PORT")
	v.BindEnv("db.user", "DB_USER")
	v.BindEnv("db.password", "DB_PASSWORD")
	v.BindEnv("db.dbname", "DB_NAME")
	v.BindEnv("db.sslmode", "DB_SSLMODE")
	v.BindEnv("jwt.secret", "JWT_SECRET_KEY")

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
