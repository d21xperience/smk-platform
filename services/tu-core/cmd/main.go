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
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"

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
	"sekolah-platform/services/tu-core/internal/storage"

	"sekolah-platform/platform/notification/channels/email"
	"sekolah-platform/platform/notification/channels/whatsapp"
	notifDispatcher "sekolah-platform/platform/notification/dispatcher"
	notifTemplates "sekolah-platform/platform/notification/templates"

	classEngine "sekolah-platform/services/tu-core/internal/domain/academic/engine"
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

	dbURL := getEnv("DATABASE_URL", "postgres://-smk_dev:-secret123@localhost:5432/-smk_platform?sslmode=disable")
	db, err := initDatabase(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database terhubung dan terverifikasi")

	sqlDB := stdlib.OpenDBFromPool(db)
	defer sqlDB.Close()
	prefService := service.NewPreferenceService(sqlDB)
	log.Println("✅ PreferenceService diinisialisasi")

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	nc, js, err := initNATS(natsURL)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke NATS: %v", err)
	}
	defer nc.Close()
	log.Println("✅ NATS JetStream terhubung")

	minioClient, err := storage.NewMinioClient()
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi MinIO: %v", err)
	}
	log.Println("✅ MinIO Storage terhubung")

	notifDispatch := setupNotificationPlatform()

	// REPOSITORIES
	studentRepo := repository.NewStudentRepository(db)
	notifRepo := notification.NewNotificationRepositoryImpl(db)
	mutationRepository := repository.NewMutationRepository(db)
	registrationRepo := repository.NewRegistrationRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	assessmentRepo := repository.NewAssessmentRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	correspondenceRepo := repository.NewCorrespondenceRepository(db)
	dispositionRepo := repository.NewDispositionRepository(db)
	signatureRepo := repository.NewSignatureRepository(db)
	classRepo := repository.NewClassRepository(db)
	statsRepo := repository.NewStatisticsRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	teacherHistoryRepo := repository.NewTeacherHistoryRepository(db)
	ptkLeaveRepo := repository.NewPtkLeaveRepository(db)
	ptkPerformanceRepo := repository.NewPtkPerformanceRepository(db) // <-- NEW
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	userRepo := repository.NewUserRepository(db)
	publicRepo := repository.NewPublicRepository(db)
	academicRepo := repository.NewAcademicRepository(db)
	schoolRepo := repository.NewSchoolRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	pembelajaranRepo := repository.NewPembelajaranRepository(db)
	kurikulumStatRepo := repository.NewKurikulumStatRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)

	// ENGINES
	studentEngine := studentEngine.NewStudentEngine()
	mutEngine := mutationEngine.NewMutationEngine()
	regEngine := registrationEngine.NewRegistrationEngine()
	attEngine := attendEngine.NewAttendanceEngine()
	assEngine := assessmentEngine.NewAssessmentEngine()
	finEngine := financeEngine.NewFinanceEngine()
	corrEngine := correspondenceEngine.NewCorrespondenceEngine()
	classEngine := classEngine.NewClassEngine()

	eventPublisher := &natsEventPublisher{nc: nc, js: js}

	// SERVICES
	studentService := service.NewStudentService(studentRepo, studentEngine, eventPublisher)
	mutationService := service.NewMutationService(mutationRepository, studentRepo, mutEngine, eventPublisher)
	registrationService := service.NewRegistrationService(registrationRepo, studentRepo, regEngine, eventPublisher)
	attendanceService := service.NewAttendanceService(attendanceRepo, studentRepo, attEngine, eventPublisher)
	assessmentService := service.NewAssessmentService(assessmentRepo, studentRepo, assEngine, eventPublisher)
	financeService := service.NewFinanceService(financeRepo, studentRepo, finEngine, eventPublisher)
	signatureService := service.NewSignatureService(signatureRepo)
	correspondenceService := service.NewCorrespondenceService(correspondenceRepo, signatureService, corrEngine, eventPublisher)
	dispositionService := service.NewDispositionService(dispositionRepo)
	classService := service.NewClassService(classRepo, classEngine, eventPublisher)
	statsService := service.NewStatisticsService(statsRepo)
	teacherService := service.NewTeacherService(teacherRepo)
	teacherHistoryService := service.NewTeacherHistoryService(teacherHistoryRepo)
	ptkLeaveService := service.NewPtkLeaveService(ptkLeaveRepo)
	ptkPerformanceService := service.NewPtkPerformanceService(ptkPerformanceRepo) // <-- NEW
	kurikulumService := service.NewKurikulumService(subjectRepo, pembelajaranRepo, kurikulumStatRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, db)
	enrollmentService := service.NewEnrollmentService(db, enrollmentRepo, classRepo, studentRepo, studentEngine, eventPublisher)
	authService := service.NewAuthService(userRepo, academicRepo)
	publicService := service.NewPublicService(publicRepo)
	userService := service.NewUserService(userRepo)
	schoolService := service.NewSchoolService(schoolRepo)

	setupBackgroundWorkers(ctx, db, js, notifDispatch, notifRepo)

	sched := setupScheduler(db, js)
	sched.Start()
	defer sched.Stop()

	grpcConfig := server.Config{Port: getEnvAsInt("GRPC_PORT", 50051)}
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

	httpConfig := server.HTTPServerConfig{Port: getEnvAsInt("HTTP_PORT", 8080)}
	healthHandler := health.NewHealthHandler(db, nc, js)
	httpServer := server.NewHTTPServer(httpConfig, healthHandler)

	restConfig := server.RESTServerConfig{Port: getEnvAsInt("REST_PORT", 8081)}
	scheduleHandler := rest.NewScheduleHandler(classRepo, academicRepo, db)

	restHandler := rest.NewRouter(
		studentService,
		mutationService,
		registrationService,
		attendanceService,
		assessmentService,
		financeService,
		correspondenceService,
		classService,
		enrollmentService,
		statsService,
		teacherService,
		kurikulumService,
		scheduleService,
		prefService,
		authService,
		publicService,
		userService,
		schoolService,
		scheduleHandler,
		dispositionService,
		signatureService,
		minioClient,
		teacherHistoryService, // <-- NEW
		ptkLeaveService, // <-- NEW
		ptkPerformanceService, // <-- NEW
	)
	restServer := server.NewRESTServer(restConfig, restHandler)

	appServer := server.NewAppServer(grpcServer, httpServer, restServer)
	appServer.Start()

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

func initDatabase(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil { return nil, fmt.Errorf("gagal parse config database: %w", err) }
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2
	dbConfig.MaxConnLifetime = 30 * time.Minute
	dbConfig.MaxConnIdleTime = 5 * time.Minute
	db, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil { return nil, fmt.Errorf("gagal membuat pool: %w", err) }
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()
	if err := db.Ping(pingCtx); err != nil { db.Close(); return nil, fmt.Errorf("gagal koneksi ke database: %w", err) }
	return db, nil
}

func initNATS(natsURL string) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(natsURL, nats.Name("tu-core-service"))
	if err != nil { return nil, nil, fmt.Errorf("gagal koneksi ke NATS: %w", err) }
	js, err := nc.JetStream()
	if err != nil { nc.Close(); return nil, nil, fmt.Errorf("gagal inisialisasi JetStream: %w", err) }
	return nc, js, nil
}

func setupNotificationPlatform() *notifDispatcher.NotificationDispatcher {
	templater := notifTemplates.NewTemplateEngine(notifTemplates.NewRegistry())
	for _, tmpl := range notifTemplates.DefaultTemplates() { templater.Register(tmpl) }
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
	outboxWorker := outbox.NewOutboxWorker(db, js)
	go outboxWorker.Start(ctx)
	projManager := projection.NewProjectionManager()
	projManager.Register(projection.NewStudentProjection(db))
	log.Println("✅ Background Workers diinisialisasi")
}

func setupScheduler(db *pgxpool.Pool, js nats.JetStreamContext) *scheduler.Scheduler {
	sched := scheduler.NewScheduler()
	sched.RegisterJob("0 2 * * *", jobs.NewOverdueInvoiceJob(db, js))
	return sched
}

type natsEventPublisher struct { nc *nats.Conn; js nats.JetStreamContext }
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
	for _, e := range events { if err := p.Publish(ctx, e); err != nil { return err } }
	return nil
}
func getEnv(key, fallback string) string { if value, exists := os.LookupEnv(key); exists { return value }; return fallback }
func getEnvAsInt(key string, fallback int) int { if valueStr, exists := os.LookupEnv(key); exists { var value int; if _, err := fmt.Sscanf(valueStr, "%d", &value); err == nil { return value } }; return fallback }
type dummyLogger struct{}
func (l *dummyLogger) Info(msg string, keysAndValues ...interface{}) { log.Printf("[INFO] %s %v", msg, keysAndValues) }
func (l *dummyLogger) Error(msg string, keysAndValues ...interface{}) { log.Printf("[ERROR] %s %v", msg, keysAndValues) }
func (l *dummyLogger) Debug(msg string, keysAndValues ...interface{}) { log.Printf("[DEBUG] %s %v", msg, keysAndValues) }
