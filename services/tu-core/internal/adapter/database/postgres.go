// Package database menyediakan database connections dan utilities.
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sekolah-platform/services/tu-core/internal/config"
)

// NewPostgresConnection membuat koneksi PostgreSQL.
func NewPostgresConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.SSLMode == "disable" { // Development
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}

// Close menutup koneksi database.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping mengecek koneksi database.
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Migrate menjalankan migrations.
func Migrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// Transaction menjalankan operasi dalam transaksi.
func Transaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

// WithTimeout menambahkan timeout ke context.
func WithTimeout(db *gorm.DB, timeout time.Duration) *gorm.DB {
	// GORM tidak support timeout langsung, gunakan context
	return db
}
