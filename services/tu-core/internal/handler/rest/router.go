package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	platformctx "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

// Router mengimplementasikan http.Handler untuk REST API
type Router struct {
	studentService        *service.StudentService
	mutationService       *service.MutationService
	registrationService   *service.RegistrationService
	attendanceService     *service.AttendanceService
	assessmentService     *service.AssessmentService
	financeService        *service.FinanceService
	correspondenceService *service.CorrespondenceService
}

// NewRouter membuat instance router REST baru
func NewRouter(
	studentService *service.StudentService,
	mutationService *service.MutationService,
	registrationService *service.RegistrationService,
	attendanceService *service.AttendanceService,
	assessmentService *service.AssessmentService,
	financeService *service.FinanceService,
	correspondenceService *service.CorrespondenceService,
) http.Handler {
	r := &Router{
		studentService:        studentService,
		mutationService:       mutationService,
		registrationService:   registrationService,
		attendanceService:     attendanceService,
		assessmentService:     assessmentService,
		financeService:        financeService,
		correspondenceService: correspondenceService,
	}

	mux := http.NewServeMux()

	// === STUDENT ROUTES ===
	mux.HandleFunc("/api/v1/students", r.handleStudents)
	mux.HandleFunc("/api/v1/students/", r.handleStudentByID)

	// === STUBS FOR OTHER DOMAINS ===
	mux.HandleFunc("/api/v1/mutations", r.stubHandler)
	mux.HandleFunc("/api/v1/registrations", r.stubHandler)
	mux.HandleFunc("/api/v1/attendance", r.stubHandler)
	mux.HandleFunc("/api/v1/assessments", r.stubHandler)
	mux.HandleFunc("/api/v1/finance", r.stubHandler)
	mux.HandleFunc("/api/v1/correspondence", r.stubHandler)

	return mux
}

// enrichContext mengekstrak header dan menyuntikkan OperationalContext ke dalam context
func (r *Router) enrichContext(req *http.Request) (context.Context, error) {
	schoolID := req.Header.Get("x-school-id")
	academicYear := req.Header.Get("x-academic-year")
	semesterStr := req.Header.Get("x-semester")
	academicPeriodID := req.Header.Get("x-academic-period")
	userID := req.Header.Get("x-user-id")
	role := req.Header.Get("x-user-role")

	// Validasi minimal sesuai validateOperationalContext di middleware
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

	// Suntikkan ke context agar GetOperationalContext(ctx) di service layer berhasil
	ctx := middleware.WithOperationalContext(req.Context(), opCtx)
	return ctx, nil
}

func (r *Router) handleStudents(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers (x-school-id, x-academic-year, x-semester, x-user-id, x-user-role)"}`, http.StatusBadRequest)
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
		log.Printf("❌ [REST] Invalid request body: %v", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	log.Printf("📝 [REST] Creating student - School: %s, Year: %s, NIS: %s", opCtx.SchoolID, opCtx.AcademicYear, reqBody.NIS)

	// Panggil service dengan signature yang sesuai dengan existing code Anda
	student, err := r.studentService.CreateStudent(
		req.Context(),
		opCtx.SchoolID,
		opCtx.AcademicPeriodID,
		opCtx.Semester,
		opCtx.UserID,
		opCtx.Role,
		reqBody,
	)
	if err != nil {
		log.Printf("❌ [REST] Error creating student: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [REST] Student created successfully with ID: %s", student.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func (r *Router) listStudents(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())

	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	offset, _ := strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	filter := service.StudentFilterRequest{
		Limit:  limit,
		Offset: offset,
	}

	students, err := r.studentService.ListStudents(
		req.Context(),
		opCtx.SchoolID,
		opCtx.AcademicPeriodID,
		filter,
	)
	if err != nil {
		log.Printf("❌ [REST] Error listing students: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (r *Router) handleStudentByID(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Endpoint ini akan segera diimplementasikan"})
}

func (r *Router) stubHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Domain endpoint ini dalam tahap pengembangan"})
}
