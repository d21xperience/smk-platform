package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"sekolah-platform/platform/message-broker/contract"
	"sekolah-platform/platform/message-broker/nats"
	"sekolah-platform/services/edge-bff/internal/client"
	"sekolah-platform/services/edge-bff/internal/config"
	"sekolah-platform/services/edge-bff/internal/handler/rest"
	"sekolah-platform/services/edge-bff/internal/handler/rest/middleware"
	"sekolah-platform/services/edge-bff/internal/projection"
	"sekolah-platform/services/edge-bff/internal/repository"
	"sekolah-platform/services/edge-bff/internal/subscriber"
)

func main() {
	log.Println("🚀 Starting edge-bff service...")

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	log.Printf("✅ Config loaded (service: %s, env: %s)", cfg.App.ServiceName, cfg.App.Environment)

	// 2. Setup SQLite read model
	repo, err := repository.NewSQLiteRepository(cfg.SQLite.Path)
	if err != nil {
		log.Fatalf("❌ Failed to setup SQLite: %v", err)
	}
	log.Printf("✅ SQLite read model ready at %s", cfg.SQLite.Path)
	defer repo.Close()

	// 3. Connect to tu-core via gRPC
	grpcClient, err := client.NewGRPCClient(cfg.TuCore.Address, cfg.TuCore.Timeout)
	if err != nil {
		log.Fatalf("❌ Failed to connect to tu-core: %v", err)
	}
	log.Printf("✅ Connected to tu-core at %s", cfg.TuCore.Address)
	defer grpcClient.Close()

	// 4. Connect to NATS
	natsConn, err := nats.NewConnection(nats.Config{
		URL:            cfg.NATS.URL,
		ConnectionName: cfg.App.ServiceName,
	}, nil)
	if err != nil {
		log.Fatalf("❌ Failed to connect to NATS: %v", err)
	}
	log.Printf("✅ Connected to NATS at %s", cfg.NATS.URL)
	defer natsConn.Close()

	// 5. Setup projection manager
	projManager := projection.NewProjectionManager(repo)
	log.Printf("✅ Projection manager ready (%d handlers)", projManager.HandlerCount())

	// 6. Start event subscriber
	natsSub := nats.NewNATSSubscriber(natsConn, cfg.NATS.StreamName, cfg.NATS.DurableName)
	eventSub := subscriber.NewEventSubscriber(natsSub, projManager, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := eventSub.Start(ctx); err != nil {
		log.Fatalf("❌ Failed to start event subscriber: %v", err)
	}
	log.Println("✅ Event subscriber started")

	// 7. Setup HTTP server (Gin)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(requestLogger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": cfg.App.ServiceName,
		})
	})

	// API routes dengan middleware
	api := router.Group("/api/v1")
	api.Use(middleware.OperationalContextMiddleware())
	// api.Use(middleware.JWTMiddleware(middleware.JWTConfig{
	// 	Secret: cfg.JWT.Secret,
	// 	Issuer: cfg.JWT.Issuer,
	// }))

	// Register handlers
	studentHandler := rest.NewStudentHandler(repo)
	studentHandler.RegisterRoutes(api)

	teacherHandler := rest.NewTeacherHandler(repo)
	teacherHandler.RegisterRoutes(api)

	// 8. Start HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	go func() {
		log.Printf("🚀 HTTP server listening on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	// 9. Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("🛑 Received signal %v, shutting down gracefully...", sig)

	// 10. Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	log.Println("⏳ Stopping HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  HTTP server shutdown error: %v", err)
	}
	log.Println("✅ HTTP server stopped")

	log.Println("⏳ Stopping event subscriber...")
	eventSub.Stop()
	log.Println("✅ Event subscriber stopped")

	log.Println("⏳ Closing NATS connection...")
	natsConn.Close()
	log.Println("✅ NATS connection closed")

	log.Println("⏳ Closing tu-core connection...")
	grpcClient.Close()
	log.Println("✅ tu-core connection closed")

	log.Println("⏳ Closing SQLite...")
	repo.Close()
	log.Println("✅ SQLite closed")

	log.Println("👋 edge-bff service stopped gracefully")
}

// requestLogger adalah middleware untuk logging request.
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("[%s] %s %s - %d (%v)",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
		)
	}
}

// Placeholder untuk contract.StreamConfig
var _ contract.StreamConfig
