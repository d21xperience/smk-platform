package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	"sekolah-platform/services/tu-core/internal/storage"
	"sekolah-platform/platform/notification/channels/email"
	"sekolah-platform/platform/notification/channels/whatsapp"
	notifDispatcher "sekolah-platform/platform/notification/dispatcher"
	notifTemplates "sekolah-platform/platform/notification/templates"
)

type Infrastructure struct {
	DB             *pgxpool.Pool
	NC             *nats.Conn
	JS             nats.JetStreamContext
	MinIO          *storage.MinioClient
	NotifDispatch  *notifDispatcher.NotificationDispatcher
	EventPublisher *natsEventPublisher
}

func (i *Infrastructure) Cleanup() {
	i.DB.Close()
	i.NC.Close()
}

func initInfrastructure(ctx context.Context) *Infrastructure {
	dbURL := getEnv("DATABASE_URL", "postgres://-smk_dev:-secret123@localhost:5432/-smk_platform?sslmode=disable")
	db, err := initDatabase(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi database: %v", err)
	}

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	nc, js, err := initNATS(natsURL)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi NATS: %v", err)
	}

	minioClient, err := storage.NewMinioClient()
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi MinIO: %v", err)
	}

	notifDispatch := setupNotificationPlatform()
	eventPublisher := &natsEventPublisher{nc: nc, js: js}

	log.Println("✅ Infrastructure terhubung")
	return &Infrastructure{
		DB: db, NC: nc, JS: js,
		MinIO: minioClient, NotifDispatch: notifDispatch, EventPublisher: eventPublisher,
	}
}

func initDatabase(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("gagal parse config database: %w", err)
	}
	dbConfig.MaxConns, dbConfig.MinConns = 10, 2
	dbConfig.MaxConnLifetime, dbConfig.MaxConnIdleTime = 30*time.Minute, 5*time.Minute
	db, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pool: %w", err)
	}
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()
	if err := db.Ping(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal koneksi ke database: %w", err)
	}
	return db, nil
}

func initNATS(natsURL string) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(natsURL, nats.Name("tu-core-service"))
	if err != nil {
		return nil, nil, fmt.Errorf("gagal koneksi ke NATS: %w", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, nil, fmt.Errorf("gagal inisialisasi JetStream: %w", err)
	}
	return nc, js, nil
}

func setupNotificationPlatform() *notifDispatcher.NotificationDispatcher {
	templater := notifTemplates.NewTemplateEngine(notifTemplates.NewRegistry())
	for _, tmpl := range notifTemplates.DefaultTemplates() {
		templater.Register(tmpl)
	}
	dispatch := notifDispatcher.NewNotificationDispatcher(nil)
	dispatch.Register("whatsapp", whatsapp.NewWhatsAppChannel(whatsapp.Config{
		GatewayURL: getEnv("WA_GATEWAY_URL", "https://api.fonnte.com/send"),
		APIKey:     getEnv("WA_API_KEY", "dummy_key"),
		Timeout:    10 * time.Second,
	}, templater))
	dispatch.Register("email", email.NewEmailChannel(email.Config{
		Host:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:      getEnv("SMTP_PORT", "587"),
		Username:  getEnv("SMTP_USER", "noreply@sekolah.sch.id"),
		Password:  getEnv("SMTP_PASS", "dummy_pass"),
		FromName:  "SDP Sekolah",
		FromEmail: "noreply@sekolah.sch.id",
	}, templater))
	return dispatch
}
