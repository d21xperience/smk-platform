package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	// Domain & Services

	"sekolah-platform/services/tu-core/internal/handler/rest"
	"sekolah-platform/services/tu-core/internal/health"
	"sekolah-platform/services/tu-core/internal/notification"
	notifHandlers "sekolah-platform/services/tu-core/internal/notification/handlers"
	"sekolah-platform/services/tu-core/internal/outbox"
	"sekolah-platform/services/tu-core/internal/projection"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/scheduler"
	"sekolah-platform/services/tu-core/internal/scheduler/jobs"
	"sekolah-platform/services/tu-core/internal/server"
	"sekolah-platform/services/tu-core/internal/service"

	// Platform Runtime
	"sekolah-platform/platform/notification/channels/email"
	"sekolah-platform/platform/notification/channels/whatsapp"
	notifDispatcher "sekolah-platform/platform/notification/dispatcher"
	notifTemplates "sekolah-platform/platform/notification/templates"

	// Domain Engines (GUNAKAN ALIAS AGAR TIDAK BENTROK)
	assessmentEngine "sekolah-platform/services/tu-core/internal/domain/assessment/engine"
	attendEngine "sekolah-platform/services/tu-core/internal/domain/attendance/engine"
	correspondenceEngine "sekolah-platform/services/tu-core/internal/domain/correspondence/engine"
	financeEngine "sekolah-platform/services/tu-core/internal/domain/finance/engine"
	mutationEngine "sekolah-platform/services/tu-core/internal/domain/mutation/engine"
	registrationEngine "sekolah-platform/services/tu-core/internal/domain/registration/engine"
	studentEngine "sekolah-platform/services/tu-core/internal/domain/student/engine"
)

func main() {
	log.Println("🚀 Memulai School Digital Platform (tu-core)...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// =========================================================================
	// 1. DATABASE CONNECTION (dengan validasi)
	// =========================================================================
	// DIPERBAIKI: Default URL match dengan .env
	dbURL := getEnv("DATABASE_URL", "postgres://-smk_dev:-secret123@localhost:5432/-smk_platform?sslmode=disable")

	db, err := initDatabase(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database terhubung dan terverifikasi")

	// =========================================================================
	// 2. NATS CONNECTION
	// =========================================================================
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	nc, js, err := initNATS(natsURL)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke NATS: %v", err)
	}
	defer nc.Close()
	log.Println("✅ NATS JetStream terhubung")

	// =========================================================================
	// 3. PLATFORM NOTIFICATION
	// =========================================================================
	notifDispatch := setupNotificationPlatform()

	// =========================================================================
	// 4. REPOSITORIES
	// =========================================================================
	studentRepo := repository.NewStudentRepository(db)
	notifRepo := notification.NewNotificationRepositoryImpl(db)
	mutationRepository := repository.NewMutationRepository(db)
	registrationRepo := repository.NewRegistrationRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	assessmentRepo := repository.NewAssessmentRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	correspondenceRepo := repository.NewCorrespondenceRepository(db)
	// =========================================================================
	// 5. DOMAIN ENGINES
	// =========================================================================
	studentEngine := studentEngine.NewStudentEngine()
	mutEngine := mutationEngine.NewMutationEngine()
	regEngine := registrationEngine.NewRegistrationEngine() // <-- TAMBAHKAN
	attEngine := attendEngine.NewAttendanceEngine()
	assEngine := assessmentEngine.NewAssessmentEngine()
	finEngine := financeEngine.NewFinanceEngine()
	corrEngine := correspondenceEngine.NewCorrespondenceEngine()
	// =========================================================================
	// 6. EVENT PUBLISHER
	// =========================================================================
	eventPublisher := &natsEventPublisher{nc: nc, js: js}

	// =========================================================================
	// 7. APPLICATION SERVICES
	// =========================================================================
	studentService := service.NewStudentService(studentRepo, studentEngine, eventPublisher)
	mutationService := service.NewMutationService(mutationRepository, studentRepo, mutEngine, eventPublisher)
	registrationService := service.NewRegistrationService(registrationRepo, studentRepo, regEngine, eventPublisher) // <-- TAMBAHKAN
	attendanceService := service.NewAttendanceService(attendanceRepo, studentRepo, attEngine, eventPublisher)
	assessmentService := service.NewAssessmentService(assessmentRepo, studentRepo, assEngine, eventPublisher)
	financeService := service.NewFinanceService(financeRepo, studentRepo, finEngine, eventPublisher)
	correspondenceService := service.NewCorrespondenceService(correspondenceRepo, corrEngine, eventPublisher)
	// =========================================================================
	// 8. BACKGROUND WORKERS
	// =========================================================================
	setupBackgroundWorkers(ctx, db, js, notifDispatch, notifRepo)

	// =========================================================================
	// 9. SCHEDULER
	// =========================================================================
	sched := setupScheduler(db, js)
	sched.Start()
	defer sched.Stop()

	// =========================================================================
	// 10. SERVERS (gRPC + HTTP)
	// =========================================================================
	// --- 1. gRPC Server ---
	grpcConfig := server.Config{Port: getEnvAsInt("GRPC_PORT", 50051)}
	// grpcServer, err := server.NewGRPCServer(grpcConfig, studentService, mutationSvc, &dummyLogger{})
	grpcServer, err := server.NewGRPCServer(
		grpcConfig,
		studentService,
		mutationService,
		registrationService,
		attendanceService,
		assessmentService,
		financeService,
		correspondenceService,
		&dummyLogger{},
	)
	if err != nil {
		log.Fatalf("❌ Gagal membuat gRPC Server: %v", err)
	}

	// --- 2. HTTP Server (Health/Metrics) ---
	httpConfig := server.HTTPServerConfig{Port: getEnvAsInt("HTTP_PORT", 8080)}
	healthHandler := health.NewHealthHandler(db, nc, js)
	httpServer := server.NewHTTPServer(httpConfig, healthHandler)

	// --- 3. REST Server (NEW: Untuk Frontend Integration) ---
	restConfig := server.RESTServerConfig{Port: getEnvAsInt("REST_PORT", 8081)} // Port 8081 agar tidak bentrok dengan 8080
	restHandler := rest.NewRouter(
		studentService,
		mutationService,
		registrationService,
		attendanceService,
		assessmentService,
		financeService,
		correspondenceService,
	)
	restServer := server.NewRESTServer(restConfig, restHandler)

	// --- 4. AppServer (Mengelola ketiganya) ---
	appServer := server.NewAppServer(grpcServer, httpServer, restServer) // <-- UPDATE: Tambahkan restServer
	appServer.Start()

	// =========================================================================
	// 11. GRACEFUL SHUTDOWN
	// =========================================================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n⏳ Sinyal shutdown diterima. Memulai graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	appServer.Stop(shutdownCtx)
	nc.Drain()

	log.Println("✅ tu-core berhasil dihentikan dengan aman.")
}

// =========================================================================
// INITIALIZATION FUNCTIONS
// =========================================================================

// initDatabase membuat koneksi database dan MEMVALIDASInya dengan Ping.
func initDatabase(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("gagal parse config database: %w", err)
	}

	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2
	dbConfig.MaxConnLifetime = 30 * time.Minute
	dbConfig.MaxConnIdleTime = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pool: %w", err)
	}

	// VALIDASI KONEKSI dengan Ping
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	if err := db.Ping(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal koneksi ke database (ping failed): %w\n\n"+
			"💡 KEMUNGKINAN PENYEBAB:\n"+
			"   1. PostgreSQL belum berjalan\n"+
			"   2. DATABASE_URL environment variable tidak di-set atau salah\n"+
			"   3. Username/password database salah\n"+
			"   4. Database belum dibuat\n\n"+
			"💡 SOLUSI:\n"+
			"   - Jalankan: make infra-up\n"+
			"   - Jalankan: make db-seed\n"+
			"   - Periksa file .env", err)
	}

	return db, nil
}

// initNATS membuat koneksi NATS dan menginisialisasi JetStream.
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

// =========================================================================
// SETUP FUNCTIONS
// =========================================================================

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

func setupBackgroundWorkers(ctx context.Context, db *pgxpool.Pool, js nats.JetStreamContext, _ *notifDispatcher.NotificationDispatcher, _ notifHandlers.InvoiceRepository) {
	// 1. Outbox Worker
	outboxWorker := outbox.NewOutboxWorker(db, js)
	go outboxWorker.Start(ctx)

	// 2. Projection Manager
	projManager := projection.NewProjectionManager()
	projManager.Register(projection.NewStudentProjection(db))

	log.Println("✅ Background Workers diinisialisasi")
}

func setupScheduler(db *pgxpool.Pool, js nats.JetStreamContext) *scheduler.Scheduler {
	sched := scheduler.NewScheduler()
	sched.RegisterJob("0 2 * * *", jobs.NewOverdueInvoiceJob(db, js))
	return sched
}

// =========================================================================
// ADAPTERS & HELPERS
// =========================================================================

type natsEventPublisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func (p *natsEventPublisher) Publish(ctx context.Context, event service.DomainEvent) error {
	payload := []byte(fmt.Sprintf("%+v", event.Payload))
	msg := &nats.Msg{
		Subject: fmt.Sprintf("sdp.event.%s.%s", event.AggregateType, event.EventType),
		Data:    payload,
		Header:  make(nats.Header),
	}
	msg.Header.Set("X-Correlation-ID", event.CorrelationID)
	msg.Header.Set("X-Event-Type", event.EventType)
	_, err := p.js.PublishMsg(msg)
	return err
}

func (p *natsEventPublisher) PublishBatch(ctx context.Context, events []service.DomainEvent) error {
	for _, e := range events {
		if err := p.Publish(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		var value int
		if _, err := fmt.Sscanf(valueStr, "%d", &value); err == nil {
			return value
		}
	}
	return fallback
}

type dummyLogger struct{}

func (l *dummyLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *dummyLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, keysAndValues)
}

func (l *dummyLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, keysAndValues)
}
