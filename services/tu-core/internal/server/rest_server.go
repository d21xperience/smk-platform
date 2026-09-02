package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// RESTServerConfig konfigurasi untuk REST API server.
type RESTServerConfig struct {
	Port int
}

// RESTServer membungkus http.Server untuk graceful shutdown.
type RESTServer struct {
	server *http.Server
}

// NewRESTServer membuat instance REST server baru.
func NewRESTServer(config RESTServerConfig, handler http.Handler) *RESTServer {
	return &RESTServer{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Port),
			Handler:      handler,
			ReadTimeout:  10 * time.Second, // Mencegah Slowloris attack
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

// Start memulai REST server (biasanya dijalankan dalam goroutine).
func (s *RESTServer) Start() error {
	return s.server.ListenAndServe()
}

// Stop menghentikan REST server dengan graceful shutdown.
func (s *RESTServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
