package nats

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Connection mengelola koneksi ke NATS server.
type Connection struct {
	config Config
	conn   *nats.Conn
	js     nats.JetStreamContext
	mu     sync.RWMutex
	logger Logger
}

// Logger adalah interface untuk logging.
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	log.Printf("[NATS][INFO] "+msg, args...)
}
func (l *defaultLogger) Error(msg string, args ...interface{}) {
	log.Printf("[NATS][ERROR] "+msg, args...)
}
func (l *defaultLogger) Debug(msg string, args ...interface{}) {}

// NewConnection membuat Connection baru.
func NewConnection(config Config, logger Logger) (*Connection, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if logger == nil {
		logger = &defaultLogger{}
	}

	c := &Connection{
		config: config,
		logger: logger,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

// connect melakukan koneksi ke NATS server.
func (c *Connection) connect() error {
	opts := []nats.Option{
		nats.Name(c.config.ConnectionName),
		nats.ReconnectWait(c.config.ReconnectWait),
		nats.MaxReconnects(c.config.MaxReconnects),
		nats.Timeout(c.config.Timeout),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			c.logger.Error("Disconnected from NATS: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			c.logger.Info("Reconnected to NATS at %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			c.logger.Info("NATS connection closed")
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			c.logger.Error("NATS error: %v", err)
		}),
	}

	// Add authentication if configured
	if c.config.CredentialsFile != "" {
		opts = append(opts, nats.UserCredentials(c.config.CredentialsFile))
	}
	if c.config.Token != "" {
		opts = append(opts, nats.Token(c.config.Token))
	}

	conn, err := nats.Connect(c.config.URL, opts...)
	if err != nil {
		return &ConnectionError{URL: c.config.URL, Cause: err}
	}

	// Get JetStream context
	var js nats.JetStreamContext
	if c.config.JetStreamDomain != "" {
		js, err = conn.JetStream(nats.Domain(c.config.JetStreamDomain))
	} else {
		js, err = conn.JetStream()
	}
	if err != nil {
		conn.Close()
		return fmt.Errorf("gagal mendapatkan JetStream context: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.js = js
	c.mu.Unlock()

	c.logger.Info("Connected to NATS at %s", conn.ConnectedUrl())
	return nil
}

// Conn mengembalikan underlying NATS connection.
func (c *Connection) Conn() *nats.Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

// JetStream mengembalikan JetStream context.
func (c *Connection) JetStream() nats.JetStreamContext {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.js
}

// IsConnected mengecek apakah connection masih aktif.
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && c.conn.IsConnected()
}

// Close menutup connection.
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.js = nil
	}
	return nil
}

// WaitForConnection menunggu sampai connection siap.
func (c *Connection) WaitForConnection(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if c.IsConnected() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errNotConnected
}
