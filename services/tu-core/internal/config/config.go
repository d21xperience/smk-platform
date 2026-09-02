// Package config menyediakan konfigurasi untuk tu-core service.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config adalah konfigurasi utama tu-core.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	NATS     NATSConfig
	Server   ServerConfig
}

// AppConfig adalah konfigurasi aplikasi.
type AppConfig struct {
	Name        string
	Environment string // development, staging, production
	LogLevel    string
}

// DatabaseConfig adalah konfigurasi database.
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// NATSConfig adalah konfigurasi NATS.
type NATSConfig struct {
	URL            string
	ConnectionName string
	StreamName     string
}

// ServerConfig adalah konfigurasi gRPC server.
type ServerConfig struct {
	Port    int
	Timeout time.Duration
}

// Load memuat konfigurasi dari environment variables.
func Load() (*Config, error) {
	// Load .env file jika ada (untuk development)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "tu-core"),
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnvInt("DB_PORT", 5432),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", ""),
			DBName:       getEnv("DB_NAME", "tu_core"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  time.Duration(getEnvInt("DB_MAX_LIFETIME_MINUTES", 5)) * time.Minute,
		},
		NATS: NATSConfig{
			URL:            getEnv("NATS_URL", "nats://localhost:4222"),
			ConnectionName: getEnv("NATS_CONNECTION_NAME", "tu-core"),
			StreamName:     getEnv("NATS_STREAM_NAME", "SDP_EVENTS"),
		},
		Server: ServerConfig{
			Port:    getEnvInt("SERVER_PORT", 50051),
			Timeout: time.Duration(getEnvInt("SERVER_TIMEOUT_SECONDS", 30)) * time.Second,
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate memvalidasi konfigurasi.
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.NATS.URL == "" {
		return fmt.Errorf("NATS_URL is required")
	}
	if c.Server.Port <= 0 {
		return fmt.Errorf("SERVER_PORT must be > 0")
	}
	return nil
}

// DSN mengembalikan Data Source Name untuk PostgreSQL.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// IsProduction mengecek apakah environment adalah production.
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment mengecek apakah environment adalah development.
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
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
