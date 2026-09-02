package server

import (
	"context"
	"log"
	"net/http"
	"sync"
)

// AppServer mengelola lifecycle dari semua server (gRPC, HTTP Health, dan REST API).
type AppServer struct {
	grpcServer *GRPCServer
	httpServer *HTTPServer
	restServer *RESTServer // <-- TAMBAHKAN
	wg         sync.WaitGroup
}

// NewAppServer membuat instance AppServer.
func NewAppServer(grpc *GRPCServer, http *HTTPServer, rest *RESTServer) *AppServer { // <-- UPDATE SIGNATURE
	return &AppServer{
		grpcServer: grpc,
		httpServer: http,
		restServer: rest, // <-- TAMBAHKAN
	}
}

// Start menjalankan semua server secara konkuren.
func (s *AppServer) Start() {
	s.wg.Add(3) // <-- UBAH dari 2 menjadi 3

	// Start gRPC Server (Blocking)
	go func() {
		defer s.wg.Done()
		log.Printf("🚀 gRPC Server berjalan di port %d", s.grpcServer.Port())
		if err := s.grpcServer.Start(); err != nil {
			log.Printf("❌ gRPC Server error: %v", err)
		}
	}()

	// Start HTTP Server (Health/Metrics) (Blocking)
	go func() {
		defer s.wg.Done()
		log.Printf("🌐 HTTP Server (Health/Metrics) berjalan di port %s", s.httpServer.server.Addr)
		if err := s.httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP Server error: %v", err)
		}
	}()

	// Start REST Server (API) (Blocking) <-- TAMBAHKAN
	go func() {
		defer s.wg.Done()
		log.Printf("🌐 REST Server (API) berjalan di port %s", s.restServer.server.Addr)
		if err := s.restServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ REST Server error: %v", err)
		}
	}()
}

// Stop melakukan graceful shutdown pada semua server.
func (s *AppServer) Stop(ctx context.Context) {
	log.Println("⏳ Memulai graceful shutdown server...")

	// 1. Stop HTTP server terlebih dahulu (lebih cepat)
	if err := s.httpServer.Stop(ctx); err != nil {
		log.Printf("⚠️ Error saat shutdown HTTP server: %v", err)
	}

	// 2. Stop REST server <-- TAMBAHKAN
	if err := s.restServer.Stop(ctx); err != nil {
		log.Printf("⚠️ Error saat shutdown REST server: %v", err)
	}

	// 3. Stop gRPC server (menunggu request yang sedang berjalan selesai)
	s.grpcServer.Stop()

	// Tunggu semua goroutine server selesai
	s.wg.Wait()
	log.Println("✅ Semua server berhasil dihentikan.")
}
