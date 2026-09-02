// Package nats menyediakan NATS JetStream implementation dari Message Broker.
package nats

import (
	"time"
)

// Config adalah konfigurasi untuk NATS connection.
type Config struct {
	// URL adalah NATS server URL (e.g. "nats://localhost:4222").
	URL string

	// ConnectionName adalah nama connection (untuk monitoring).
	ConnectionName string

	// ReconnectWait adalah durasi tunggu sebelum reconnect.
	ReconnectWait time.Duration

	// MaxReconnects adalah maksimum reconnect attempts (-1 = unlimited).
	MaxReconnects int

	// Timeout adalah timeout untuk connection dan operations.
	Timeout time.Duration

	// JetStream_domain adalah domain JetStream (untuk multi-tenant).
	JetStreamDomain string

	// Credentials untuk authentication (optional).
	CredentialsFile string

	// Token untuk token-based auth (optional).
	Token string
}

// DefaultConfig adalah konfigurasi default.
var DefaultConfig = Config{
	URL:             "nats://localhost:4222",
	ConnectionName:  "sdp-platform",
	ReconnectWait:   2 * time.Second,
	MaxReconnects:   -1, // unlimited
	Timeout:         5 * time.Second,
	JetStreamDomain: "",
}

// Validate memvalidasi konfigurasi.
func (c Config) Validate() error {
	if c.URL == "" {
		return errEmptyURL
	}
	if c.Timeout <= 0 {
		return errInvalidTimeout
	}
	return nil
}

// WithURL mengembalikan config dengan URL baru.
func (c Config) WithURL(url string) Config {
	c.URL = url
	return c
}

// WithConnectionName mengembalikan config dengan connection name baru.
func (c Config) WithConnectionName(name string) Config {
	c.ConnectionName = name
	return c
}
