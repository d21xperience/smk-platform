package rest

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	platformctx "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	"sekolah-platform/services/tu-core/internal/storage"
)

type Router struct {
	studentService        *service.StudentService
	mutationService       *service.MutationService
	registrationService   *service.RegistrationService
	attendanceService     *service.AttendanceService
	assessmentService     *service.AssessmentService
	financeService        *service.FinanceService
	correspondenceService *service.CorrespondenceService
	dispositionService    *service.DispositionService
	signatureService      *service.SignatureService
	classService          *service.ClassService
	enrollmentService     *service.EnrollmentService
	statisticsService     *service.StatisticsService
	teacherService        *service.TeacherService
	teacherHistoryService *service.TeacherHistoryService
	ptkLeaveService       *service.PtkLeaveService // <-- NEW
	kurikulumService      *service.KurikulumService
	scheduleService       service.ScheduleService
	contextHandler        *ContextHandler
	authHandler           *AuthHandler
	publicHandler         *PublicHandler
	userHandler           *UserHandler
	schoolHandler         *SchoolHandler
	attendanceHandler     *AttendanceHandler
	scheduleHandler       *ScheduleHandler
	correspondenceHandler *CorrespondenceHandler
	dispositionHandler    *DispositionHandler
	signatureHandler      *SignatureHandler
	minioClient           *storage.MinioClient
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
	prefService *service.PreferenceService,
	authService *service.AuthService,
	publicService *service.PublicService,
	userService *service.UserService,
	schoolService *service.SchoolService,
	scheduleHandler *ScheduleHandler,
	dispositionService *service.DispositionService,
	signatureService *service.SignatureService,
	minioClient *storage.MinioClient,
	teacherHistoryService *service.TeacherHistoryService,
	ptkLeaveService *service.PtkLeaveService, // <-- NEW
) http.Handler {
	r := &Router{
		studentService:        studentService,
		mutationService:       mutationService,
		registrationService:   registrationService,
		attendanceService:     attendanceService,
		assessmentService:     assessmentService,
		financeService:        financeService,
		correspondenceService: correspondenceService,
		dispositionService:    dispositionService,
		signatureService:      signatureService,
		classService:          classService,
		enrollmentService:     enrollmentService,
		statisticsService:     statisticsService,
		teacherService:        teacherService,
		teacherHistoryService: teacherHistoryService,
		ptkLeaveService:       ptkLeaveService, // <-- NEW
		kurikulumService:      kurikulumService,
		scheduleService:       scheduleService,
		contextHandler:        NewContextHandler(prefService),
		authHandler:           NewAuthHandler(authService),
		publicHandler:         NewPublicHandler(publicService),
		userHandler:           NewUserHandler(userService),
		schoolHandler:         NewSchoolHandler(schoolService),
		attendanceHandler:     NewAttendanceHandler(attendanceService),
		scheduleHandler:       scheduleHandler,
		correspondenceHandler: NewCorrespondenceHandler(correspondenceService),
		dispositionHandler:    NewDispositionHandler(dispositionService),
		signatureHandler:      NewSignatureHandler(signatureService),
		minioClient:           minioClient,
	}

	mux := http.NewServeMux()

	// Registrasi Rute Terpusat
	r.registerPtkRoutes(mux)
	r.registerStudentRoutes(mux)
	r.registerAcademicRoutes(mux)
	r.registerCorrespondenceRoutes(mux)
	r.registerSystemRoutes(mux)

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
