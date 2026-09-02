package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config adalah konfigurasi edge-bff.
type Config struct {
	App      AppConfig
	Server   ServerConfig
	TuCore   TuCoreConfig
	NATS     NATSConfig
	SQLite   SQLiteConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name        string
	Environment string
	ServiceName string // "edge-guru", "edge-website", "edge-ppdb"
}

type ServerConfig struct {
	Port    int
	Timeout time.Duration
}

type TuCoreConfig struct {
	Address string // gRPC address (e.g. "localhost:50051")
	Timeout time.Duration
}

type NATSConfig struct {
	URL           string
	DurableName   string // Consumer durable name (service-specific)
	StreamName    string
}

type SQLiteConfig struct {
	Path string // Path ke SQLite file
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	Expiration time.Duration
}

// Load memuat konfigurasi dari environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "edge-bff"),
			Environment: getEnv("APP_ENV", "development"),
			ServiceName: getEnv("SERVICE_NAME", "edge-bff"),
		},
		Server: ServerConfig{
			Port:    getEnvInt("SERVER_PORT", 8080),
			Timeout: time.Duration(getEnvInt("SERVER_TIMEOUT_SECONDS", 30)) * time.Second,
		},
		TuCore: TuCoreConfig{
			Address: getEnv("TUCORE_ADDRESS", "localhost:50051"),
			Timeout: time.Duration(getEnvInt("TUCORE_TIMEOUT_SECONDS", 30)) * time.Second,
		},
		NATS: NATSConfig{
			URL:         getEnv("NATS_URL", "nats://localhost:4222"),
			DurableName: getEnv("NATS_DURABLE_NAME", "edge-bff"),
			StreamName:  getEnv("NATS_STREAM_NAME", "SDP_EVENTS"),
		},
		SQLite: SQLiteConfig{
			Path: getEnv("SQLITE_PATH", "./data/edge-bff.db"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "dev-secret-change-in-production"),
			Issuer:     getEnv("JWT_ISSUER", "school-platform"),
			Expiration: time.Duration(getEnvInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate memvalidasi konfigurasi.
func (c *Config) Validate() error {
	if c.TuCore.Address == "" {
		return fmt.Errorf("TUCORE_ADDRESS is required")
	}
	if c.NATS.URL == "" {
		return fmt.Errorf("NATS_URL is required")
	}
	if c.SQLite.Path == "" {
		return fmt.Errorf("SQLITE_PATH is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
