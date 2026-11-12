package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("🚀 Memulai School Digital Platform (tu-core)...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Init Infrastructure
	infra := initInfrastructure(ctx)
	defer infra.Cleanup()

	// 2. Wire Dependencies
	repos := initRepositories(infra.DB)
	engines := initEngines()
	services := initServices(repos, engines, infra.MinIO, infra.EventPublisher, infra.DB)

	// 3. Background Workers & Scheduler
	setupBackgroundWorkers(ctx, infra.DB, infra.JS, infra.NotifDispatch, repos.Notification)
	sched := setupScheduler(infra.DB, infra.JS)
	sched.Start()
	defer sched.Stop()

	// 4. Start Servers
	servers := setupServers(services, repos, infra)
	servers.Start()

	// 5. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("\n🛑 Sinyal shutdown diterima. Memulai graceful shutdown...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	servers.Stop(shutdownCtx)
	infra.NC.Drain()
	log.Println("✅ tu-core berhasil dihentikan dengan aman.")
}
