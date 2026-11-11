package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	notifDispatcher "sekolah-platform/platform/notification/dispatcher"
	notifHandlers "sekolah-platform/services/tu-core/internal/notification/handlers"
	"sekolah-platform/services/tu-core/internal/outbox"
	"sekolah-platform/services/tu-core/internal/projection"
	"sekolah-platform/services/tu-core/internal/scheduler"
	"sekolah-platform/services/tu-core/internal/scheduler/jobs"
)

func setupBackgroundWorkers(ctx context.Context, db *pgxpool.Pool, js nats.JetStreamContext, _ *notifDispatcher.NotificationDispatcher, _ notifHandlers.InvoiceRepository) {
	go outbox.NewOutboxWorker(db, js).Start(ctx)
	projManager := projection.NewProjectionManager()
	projManager.Register(projection.NewStudentProjection(db))
	log.Println("✅ Background Workers diinisialisasi")
}

func setupScheduler(db *pgxpool.Pool, js nats.JetStreamContext) *scheduler.Scheduler {
	sched := scheduler.NewScheduler()
	sched.RegisterJob("0 2 * * *", jobs.NewOverdueInvoiceJob(db, js))
	return sched
}
