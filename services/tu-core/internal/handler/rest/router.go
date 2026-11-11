package rest

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	platformctx "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
	"sekolah-platform/services/tu-core/internal/storage"
)

type Router struct {
	// Akademik & Kesiswaan
	studentService      *service.StudentService
	mutationService     *service.MutationService
	registrationService *service.RegistrationService
	attendanceService   *service.AttendanceService
	assessmentService   *service.AssessmentService
	financeService      *service.FinanceService
	classService        *service.ClassService
	enrollmentService   *service.EnrollmentService
	statisticsService   *service.StatisticsService
	kurikulumService    *service.KurikulumService
	scheduleService     service.ScheduleService
	jurusanService      *service.JurusanService
	scheduleHandler     *ScheduleHandler
	attendanceHandler   *AttendanceHandler

	// Persuratan
	correspondenceService *service.CorrespondenceService
	dispositionService    *service.DispositionService
	signatureService      *service.SignatureService
	expeditionService     *service.ExpeditionService
	archiveService        *service.ArchiveService
	agendaService         *service.AgendaService
	letterTypeService     *service.LetterTypeService
	correspondenceHandler *CorrespondenceHandler
	dispositionHandler    *DispositionHandler
	signatureHandler      *SignatureHandler
	expeditionHandler     *ExpeditionHandler
	archiveHandler        *ArchiveHandler
	agendaHandler         *AgendaHandler
	// System & Auth
	contextHandler *ContextHandler
	authHandler    *AuthHandler
	publicHandler  *PublicHandler
	userHandler    *UserHandler
	schoolHandler  *SchoolHandler

	// Kepegawaian (PTK)
	teacherService           *service.TeacherService
	teacherHistoryService    *service.TeacherHistoryService
	ptkLeaveService          *service.PtkLeaveService
	ptkArchiveService        *service.PtkArchiveService
	ptkDashboardService      *service.PtkDashboardService
	bukuIndukService         *service.BukuIndukService
	laporanBerkalaRepo       repository.LaporanBerkalaRepository
	laporanAggregatorService *service.LaporanAggregatorService
	laporanBerkalaService    *service.LaporanBerkalaService
	laporanEmailService      *service.LaporanEmailService
	ptkStatisticService      *service.PtkStatisticService
	skPembagianTugasService  *service.SKPembagianTugasService
	statistikService         *service.StatistikService
	skpService               *service.SKPService
	bukuKontrolService       *service.BukuKontrolService
	usulanKenaikanService    *service.UsulanKenaikanService
	kehadiranService         *service.KehadiranService
	pelayananService         *service.PelayananService

	// Sarpras (12 Modul)
	logisticService      *service.LogisticService
	kirService           *service.KIRService
	procurementService   *service.ProcurementService
	goodsReceiptService  *service.GoodsReceiptService
	documentService      *service.DocumentService
	disposalService      *service.DisposalService
	borrowingService     *service.BorrowingService
	maintenanceService   *service.MaintenanceService
	sarprasReportService *service.SarprasReportService

	// Storage
	minioClient *storage.MinioClient
}

func NewRouter(
	studentService *service.StudentService,
	mutationService *service.MutationService,
	registrationService *service.RegistrationService,
	attendanceService *service.AttendanceService,
	assessmentService *service.AssessmentService,
	financeService *service.FinanceService,
	correspondenceService *service.CorrespondenceService,
	classService *service.ClassService,
	enrollmentService *service.EnrollmentService,
	statisticsService *service.StatisticsService,
	teacherService *service.TeacherService,
	kurikulumService *service.KurikulumService,
	scheduleService service.ScheduleService,
	jurusanService *service.JurusanService,
	prefService *service.PreferenceService,
	authService *service.AuthService,
	publicService *service.PublicService,
	userService *service.UserService,
	schoolService *service.SchoolService,
	scheduleHandler *ScheduleHandler,
	dispositionService *service.DispositionService,
	signatureService *service.SignatureService,
	expeditionService *service.ExpeditionService,
	archiveService *service.ArchiveService,
	agendaService *service.AgendaService,
	minioClient *storage.MinioClient,
	teacherHistoryService *service.TeacherHistoryService,
	ptkLeaveService *service.PtkLeaveService,
	ptkArchiveService *service.PtkArchiveService,
	ptkDashboardService *service.PtkDashboardService,
	bukuIndukService *service.BukuIndukService,
	laporanBerkalaRepo repository.LaporanBerkalaRepository,
	laporanAggregatorService *service.LaporanAggregatorService,
	laporanBerkalaService *service.LaporanBerkalaService,
	laporanEmailService *service.LaporanEmailService,
	ptkStatisticService *service.PtkStatisticService,
	skPembagianTugasService *service.SKPembagianTugasService,
	statistikService *service.StatistikService,
	skpService *service.SKPService,
	bukuKontrolService *service.BukuKontrolService,
	usulanKenaikanService *service.UsulanKenaikanService,
	kehadiranService *service.KehadiranService,
	pelayananService *service.PelayananService,
	logisticService *service.LogisticService,
	kirService *service.KIRService,
	letterTypeService *service.LetterTypeService,
	procurementService *service.ProcurementService,
	goodsReceiptService *service.GoodsReceiptService,
	documentService *service.DocumentService,
	disposalService *service.DisposalService,
	borrowingService *service.BorrowingService,
	maintenanceService *service.MaintenanceService,
	sarprasReportService *service.SarprasReportService,
) http.Handler {
	r := &Router{
		studentService:           studentService,
		mutationService:          mutationService,
		registrationService:      registrationService,
		attendanceService:        attendanceService,
		assessmentService:        assessmentService,
		financeService:           financeService,
		correspondenceService:    correspondenceService,
		dispositionService:       dispositionService,
		signatureService:         signatureService,
		expeditionService:        expeditionService,
		archiveService:           archiveService,
		agendaService:            agendaService,
		classService:             classService,
		enrollmentService:        enrollmentService,
		statisticsService:        statisticsService,
		teacherService:           teacherService,
		kurikulumService:         kurikulumService,
		scheduleService:          scheduleService,
		jurusanService:           jurusanService,
		contextHandler:           NewContextHandler(prefService),
		authHandler:              NewAuthHandler(authService),
		publicHandler:            NewPublicHandler(publicService),
		userHandler:              NewUserHandler(userService),
		schoolHandler:            NewSchoolHandler(schoolService),
		attendanceHandler:        NewAttendanceHandler(attendanceService),
		scheduleHandler:          scheduleHandler,
		correspondenceHandler:    NewCorrespondenceHandler(correspondenceService),
		dispositionHandler:       NewDispositionHandler(dispositionService),
		signatureHandler:         NewSignatureHandler(signatureService),
		expeditionHandler:        NewExpeditionHandler(expeditionService),
		archiveHandler:           NewArchiveHandler(archiveService),
		agendaHandler:            NewAgendaHandler(agendaService),
		minioClient:              minioClient,
		teacherHistoryService:    teacherHistoryService,
		ptkLeaveService:          ptkLeaveService,
		ptkArchiveService:        ptkArchiveService,
		ptkDashboardService:      ptkDashboardService,
		bukuIndukService:         bukuIndukService,
		laporanBerkalaRepo:       laporanBerkalaRepo,
		laporanAggregatorService: laporanAggregatorService,
		laporanBerkalaService:    laporanBerkalaService,
		laporanEmailService:      laporanEmailService,
		ptkStatisticService:      ptkStatisticService,
		skPembagianTugasService:  skPembagianTugasService,
		statistikService:         statistikService,
		skpService:               skpService,
		bukuKontrolService:       bukuKontrolService,
		usulanKenaikanService:    usulanKenaikanService,
		kehadiranService:         kehadiranService,
		pelayananService:         pelayananService,
		logisticService:          logisticService,
		kirService:               kirService,
		letterTypeService:        letterTypeService,
		procurementService:       procurementService,
		goodsReceiptService:      goodsReceiptService,
		documentService:          documentService,
		disposalService:          disposalService,
		borrowingService:         borrowingService,
		maintenanceService:       maintenanceService,
		sarprasReportService:     sarprasReportService,
	}

	mux := http.NewServeMux()

	// ===== RUTE KESISWAAN, AKADEMIK, KURIKULUM =====
	r.registerPtkRoutes(mux)
	r.registerStudentRoutes(mux)
	r.registerPSBRoutes(mux)
	r.registerAcademicRoutes(mux)

	// ===== RUTE PERSURATAN =====
	r.registerCorrespondenceRoutes(mux)
	r.registerArchiveRoutes(mux)
	r.registerAgendaRoutes(mux)
	r.registerLetterTypeRoutes(mux)

	// ===== RUTE SYSTEM & AUTH =====
	r.registerSystemRoutes(mux)

	// ===== RUTE STATISTIK & LAPORAN UMUM =====
	r.registerStatisticRoutes(mux)
	r.registerLaporanRoutes(mux)

	// ===== RUTE KEPEGAWAIAN (PTK) =====
	r.registerSKPRoutes(mux)
	r.registerBukuKontrolRoutes(mux)
	r.registerUsulanKenaikanRoutes(mux)
	r.registerKehadiranRoutes(mux)

	// ===== RUTE SARPRAS (12 MODUL) =====
	r.registerSarprasRoutes(mux)      // Logistik, KIR, Master Data, Mutate, History, Condition Book
	r.registerProcurementRoutes(mux)  // Pengajuan
	r.registerGoodsReceiptRoutes(mux) // Penerimaan (BAST)
	r.registerDocumentRoutes(mux)     // Dokumen Inventaris
	r.registerDisposalRoutes(mux)     // Penghapusan
	r.registerBorrowingRoutes(mux)    // Peminjaman
	r.registerMaintenanceRoutes(mux)  // Perawatan & Tiket
	r.registerReportRoutes(mux)       // Dashboard & Laporan
	// ===== RUTE STATIC FILE (EXPORT EXCEL/PDF) =====
	fs := http.FileServer(http.Dir("./storage"))
	mux.Handle("/storage/", http.StripPrefix("/storage/", fs))
	return mux
}

func (r *Router) wrapWithContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, err := r.enrichContext(req)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
			return
		}
		next(w, req.WithContext(ctx))
	}
}

func (r *Router) enrichContext(req *http.Request) (context.Context, error) {
	schoolID := req.Header.Get("x-school-id")
	academicYear := req.Header.Get("x-academic-year")
	semesterStr := req.Header.Get("x-semester")
	academicPeriodID := req.Header.Get("x-academic-period-id")
	userID := req.Header.Get("x-user-id")
	role := req.Header.Get("x-user-role")

	if schoolID == "" || userID == "" {
		return nil, http.ErrAbortHandler
	}
	if role == "" {
		role = "admin"
	}

	semester := 1
	if semesterStr != "" {
		if s, err := strconv.Atoi(semesterStr); err == nil {
			semester = s
		}
	}

	opCtx := &platformctx.OperationalContext{
		SchoolID:         schoolID,
		AcademicYear:     academicYear,
		Semester:         semester,
		AcademicPeriodID: academicPeriodID,
		UserID:           userID,
		Role:             role,
	}

	ctx := middleware.WithOperationalContext(req.Context(), opCtx)
	return ctx, nil
}

func splitPath(p string) []string {
	if p == "" {
		return []string{}
	}
	if idx := strings.Index(p, "?"); idx != -1 {
		p = p[:idx]
	}
	parts := strings.Split(p, "/")
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
