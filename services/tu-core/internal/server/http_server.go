package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HTTPServerConfig konfigurasi untuk HTTP server (health check, metrics).
type HTTPServerConfig struct {
	Port int
}

// HTTPServer membungkus http.Server untuk graceful shutdown.
type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer membuat instance HTTP server baru.
// healthHandler adalah handler yang menangani endpoint /health, /ready, /live
func NewHTTPServer(config HTTPServerConfig, healthHandler http.Handler) *HTTPServer {
	mux := http.NewServeMux()

	// Register health check endpoints
	mux.Handle("/health", healthHandler)
	mux.Handle("/ready", healthHandler)
	mux.Handle("/live", healthHandler)

	return &HTTPServer{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Port),
			Handler:      mux,
			ReadTimeout:  5 * time.Second, // Mencegah Slowloris attack
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

// Start memulai HTTP server (biasanya dijalankan dalam goroutine).
func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

// Stop menghentikan HTTP server dengan graceful shutdown.
func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
