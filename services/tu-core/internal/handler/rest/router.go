package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	platformctx "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

type Router struct {
	studentService        *service.StudentService
	mutationService       *service.MutationService
	registrationService   *service.RegistrationService
	attendanceService     *service.AttendanceService
	assessmentService     *service.AssessmentService
	financeService        *service.FinanceService
	correspondenceService *service.CorrespondenceService
	classService          *service.ClassService
	enrollmentService     *service.EnrollmentService
	statisticsService     *service.StatisticsService
	contextHandler        *ContextHandler
	authHandler           *AuthHandler
	publicHandler         *PublicHandler
	userHandler           *UserHandler
	schoolHandler         *SchoolHandler
	attendanceHandler     *AttendanceHandler
	scheduleHandler       *ScheduleHandler
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
	prefService *service.PreferenceService,
	authService *service.AuthService,
	publicService *service.PublicService,
	userService *service.UserService,
	schoolService *service.SchoolService,
	scheduleHandler *ScheduleHandler,
) http.Handler {
	r := &Router{
		studentService:        studentService,
		mutationService:       mutationService,
		registrationService:   registrationService,
		attendanceService:     attendanceService,
		assessmentService:     assessmentService,
		financeService:        financeService,
		correspondenceService: correspondenceService,
		classService:          classService,
		enrollmentService:     enrollmentService,
		statisticsService:     statisticsService,
		contextHandler:        NewContextHandler(prefService),
		authHandler:           NewAuthHandler(authService),
		publicHandler:         NewPublicHandler(publicService),
		userHandler:           NewUserHandler(userService),
		schoolHandler:         NewSchoolHandler(schoolService),
		attendanceHandler:     NewAttendanceHandler(attendanceService),
		scheduleHandler:       scheduleHandler,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/students", r.handleStudents)
	mux.HandleFunc("/api/v1/students/available", r.handleAvailableStudents)
	mux.HandleFunc("/api/v1/students/", r.handleStudentByID)
	mux.HandleFunc("/api/v1/mutations", r.stubHandler)
	mux.HandleFunc("/api/v1/registrations", r.stubHandler)
	mux.HandleFunc("/api/v1/attendance", r.stubHandler)
	mux.HandleFunc("/api/v1/assessments", r.stubHandler)
	mux.HandleFunc("/api/v1/finance", r.stubHandler)
	mux.HandleFunc("/api/v1/correspondence", r.stubHandler)
	mux.HandleFunc("/api/v1/statistics", r.handleStatistics)
	mux.HandleFunc("/api/v1/classes", r.handleClasses)
	mux.HandleFunc("/api/v1/classes/{id}", r.handleClassByID)
	mux.HandleFunc("/api/v1/classes/{id}/members", r.handleClassMembers)
	mux.HandleFunc("/api/v1/students/{id}/history", r.handleStudentHistory)

	mux.HandleFunc("/api/v1/context/periods", r.contextHandler.GetPeriods)
	mux.HandleFunc("/api/v1/users/me/preferences", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.contextHandler.GetPreference(w, req)
		case http.MethodPut:
			r.contextHandler.UpdatePreference(w, req)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})
	mux.HandleFunc("/api/v1/auth/login", r.authHandler.Login)
	mux.HandleFunc("/api/v1/public/stats", r.publicHandler.GetStats)

	// === ENDPOINT Manajemen Periode ===
	mux.HandleFunc("/api/v1/academic/periods", r.handleCreateAcademicPeriod)
	mux.HandleFunc("/api/v1/academic/periods/", r.handleAcademicPeriodsByID)

	// === ENDPOINT Manajemen Pengguna ===
	mux.HandleFunc("/api/v1/users", r.handleUsers)
	mux.HandleFunc("/api/v1/users/", r.handleUsersByID)

	// === ENDPOINT Profil Sekolah ===
	mux.HandleFunc("/api/v1/school", r.handleSchool)

	// === ENDPOINT JADWAL PELAJARAN  ===
	mux.HandleFunc("/api/v1/schedules", r.wrapWithContext(r.handleSchedules))
	mux.HandleFunc("/api/v1/schedules/options", r.wrapWithContext(r.scheduleHandler.GetOptions))

	return mux
}

// === WRAPPER JADWAL ===
func (r *Router) handleSchedules(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.scheduleHandler.ListSchedules(w, req)
	case http.MethodPost:
		r.scheduleHandler.CreateSchedule(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

// Helper Wrapper
func (r *Router) wrapWithContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, err := r.enrichContext(req)
		if err != nil {
			http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
			return
		}
		next(w, req.WithContext(ctx))
	}
}

// === WRAPPER SEKOLAH ===
func (r *Router) handleSchool(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodGet:
		r.schoolHandler.GetSchool(w, req)
	case http.MethodPut:
		r.schoolHandler.UpdateSchool(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

// === WRAPPER USER MANAGEMENT ===
func (r *Router) handleUsers(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodGet:
		r.userHandler.ListUsers(w, req)
	case http.MethodPost:
		r.userHandler.CreateUser(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleUsersByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodPut:
		r.userHandler.UpdateUser(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

// === WRAPPER PERIODE ===
func (r *Router) handleCreateAcademicPeriod(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}
	r.contextHandler.CreatePeriod(w, req.WithContext(ctx))
}

func (r *Router) handleAcademicPeriodsByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodPut:
		if strings.HasSuffix(req.URL.Path, "/activate") {
			r.contextHandler.ActivatePeriod(w, req)
			return
		}
	case http.MethodDelete:
		r.contextHandler.DeletePeriod(w, req)
		return
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
}

// === FUNGSI ASLI ===
func (r *Router) enrichContext(req *http.Request) (context.Context, error) {
	schoolID := req.Header.Get("x-school-id")
	academicYear := req.Header.Get("x-academic-year")
	semesterStr := req.Header.Get("x-semester")
	academicPeriodID := req.Header.Get("x-academic-period-id")
	userID := req.Header.Get("x-user-id")
	role := req.Header.Get("x-user-role")

	if schoolID == "" || academicYear == "" || userID == "" || role == "" {
		return nil, http.ErrAbortHandler
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

func (r *Router) handleStudents(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodPost:
		r.createStudent(w, req.WithContext(ctx))
	case http.MethodGet:
		r.listStudents(w, req.WithContext(ctx))
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (r *Router) createStudent(w http.ResponseWriter, req *http.Request) {
	var reqBody service.CreateStudentRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error": "Invalid request body: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	student, err := r.studentService.CreateStudent(
		req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID,
		opCtx.Semester, opCtx.UserID, opCtx.Role, reqBody,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, service.ErrInvalidInput) {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		} else if errors.Is(err, service.ErrDuplicateEntry) {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusConflict)
		} else {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func (r *Router) listStudents(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
	searchStr := req.URL.Query().Get("search")
	statusStr := req.URL.Query().Get("status")

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	filter := repository.StudentFilter{
		Limit:   limit,
		Offset:  offset,
		Keyword: searchStr,
		Status:  statusStr,
	}

	students, err := r.studentService.ListStudentsForUI(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID, filter)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (r *Router) handleStudentByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/students/")
	parts := splitPath(path)
	if len(parts) == 0 {
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}
	studentID := parts[0]

	opCtx := middleware.GetOperationalContext(ctx)

	switch req.Method {
	case http.MethodGet:
		student, err := r.studentService.GetStudentForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(student)

	case http.MethodPut:
		var reqBody struct {
			Name       string `json:"nama"`
			BirthPlace string `json:"tempatLahir"`
			Gender     string `json:"jenisKelamin"`
			Address    string `json:"alamat"`
			Phone      string `json:"noHp"`
			ClassName  string `json:"kelas"`
		}
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			http.Error(w, `{"error": "Invalid body"}`, http.StatusBadRequest)
			return
		}
		student, err := r.studentService.UpdateStudentForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID, reqBody)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(student)

	case http.MethodDelete:
		err := r.studentService.DeleteStudentForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (r *Router) stubHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, `{"error": "Not implemented"}`, http.StatusNotImplemented)
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

func (r *Router) handleAvailableStudents(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	if req.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	opCtx := middleware.GetOperationalContext(ctx)

	keyword := req.URL.Query().Get("keyword")
	filter := repository.StudentFilter{
		Keyword: keyword,
		Limit:   100,
		Offset:  0,
	}

	resp, err := r.studentService.ListAvailableStudentsForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, filter)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (r *Router) handleStatistics(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	if req.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	opCtx := middleware.GetOperationalContext(ctx)
	stats, err := r.statisticsService.GetDashboardStats(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
