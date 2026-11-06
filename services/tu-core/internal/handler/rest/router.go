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
	teacherService        *service.TeacherService
	kurikulumService      *service.KurikulumService
	scheduleService       service.ScheduleService
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
	teacherService *service.TeacherService,
	kurikulumService *service.KurikulumService,
	scheduleService service.ScheduleService,
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
		teacherService:        teacherService,
		kurikulumService:      kurikulumService,
		scheduleService:       scheduleService,
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
	mux.HandleFunc("/api/v1/assessments", r.stubHandler)
	mux.HandleFunc("/api/v1/finance", r.stubHandler)
	mux.HandleFunc("/api/v1/correspondence", r.stubHandler)
	mux.HandleFunc("/api/v1/teachers/{id}", r.handleTeacherByID)
	mux.HandleFunc("/api/v1/subjects", r.handleSubjects)
	mux.HandleFunc("/api/v1/subjects/{id}", r.handleSubjectByID)

	// === ENDPOINT JADWAL PELAJARAN ===
	mux.HandleFunc("/api/v1/schedules/options", r.wrapWithContext(r.scheduleHandler.GetOptions))
	mux.HandleFunc("/api/v1/schedules", r.wrapWithContext(r.handleSchedules))
	mux.HandleFunc("/api/v1/schedules/{id}", r.wrapWithContext(r.handleScheduleByID))

	// === ENDPOINT ABSENSI HARIAN ===
	// Collection: POST create session, GET list sessions
	mux.HandleFunc("/api/v1/attendance/sessions", r.wrapWithContext(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.attendanceHandler.CreateSession(w, req)
		case http.MethodGet:
			r.attendanceHandler.ListSessions(w, req)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	}))
	// Item: POST /sessions/{id}/submit
	mux.HandleFunc("/api/v1/attendance/sessions/", r.wrapWithContext(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}
		r.attendanceHandler.SubmitAttendance(w, req)
	}))
	// Recap: GET
	mux.HandleFunc("/api/v1/attendance/recap", r.wrapWithContext(r.attendanceHandler.GetRecap))

	mux.HandleFunc("/api/v1/kurikulum/statistics", r.handleKurikulumStats)
	mux.HandleFunc("/api/v1/pembelajaran", r.handlePembelajaran)
	mux.HandleFunc("/api/v1/teachers", r.handleTeachers)
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

	mux.HandleFunc("/api/v1/academic/periods", r.handleCreateAcademicPeriod)
	mux.HandleFunc("/api/v1/academic/periods/", r.handleAcademicPeriodsByID)

	mux.HandleFunc("/api/v1/users", r.handleUsers)
	mux.HandleFunc("/api/v1/users/", r.handleUsersByID)

	mux.HandleFunc("/api/v1/school", r.handleSchool)

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

// === ENRICH CONTEXT ===
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
	classID := req.URL.Query().Get("class_id") // untuk roster absensi

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
		ClassID: classID,
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

func (r *Router) handleTeachers(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	opCtx := middleware.GetOperationalContext(ctx)

	switch req.Method {
	case http.MethodGet:
		keyword := req.URL.Query().Get("keyword")
		teachers, total, err := r.teacherService.ListTeachers(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, keyword)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": teachers, "total": total})

	case http.MethodPost:
		var body service.TeacherRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, `{"error": "Invalid body"}`, http.StatusBadRequest)
			return
		}
		if err := r.teacherService.CreateTeacher(ctx, body); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Guru berhasil ditambahkan"})
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleTeacherByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context"}`, http.StatusBadRequest)
		return
	}

	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-1]
	if id == "" {
		http.Error(w, `{"error": "ID required"}`, http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodPut:
		var body service.TeacherRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, `{"error": "Invalid body"}`, http.StatusBadRequest)
			return
		}
		if err := r.teacherService.UpdateTeacher(ctx, id, body); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Guru berhasil diupdate"})

	case http.MethodDelete:
		if err := r.teacherService.DeleteTeacher(ctx, id); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Guru berhasil dihapus"})

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleSubjects(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		subjects, err := r.kurikulumService.ListSubjects(req.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"data": subjects})
	case http.MethodPost:
		var body repository.Subject
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		if err := r.kurikulumService.CreateSubject(req.Context(), &body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleSubjectByID(w http.ResponseWriter, req *http.Request) {
	// Implementasi Update & Delete Mapel (mirip handler sebelumnya)
}

func (r *Router) handlePembelajaran(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}
	opCtx := middleware.GetOperationalContext(ctx)

	switch req.Method {
	case http.MethodGet:
		classID := req.URL.Query().Get("classId")
		list, err := r.kurikulumService.ListPembelajaran(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"data": list})
	case http.MethodPost:
		var body struct {
			ClassID   string `json:"classId"`
			SubjectID string `json:"subjectId"`
			TeacherID string `json:"teacherId"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		if err := r.kurikulumService.CreatePembelajaran(ctx, body.ClassID, body.SubjectID, body.TeacherID, opCtx.AcademicPeriodID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Pembelajaran ditambahkan"})
	}
}

func (r *Router) handleKurikulumStats(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}
	opCtx := middleware.GetOperationalContext(ctx)

	stats, err := r.kurikulumService.GetDashboardStats(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func (r *Router) handleScheduleByID(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]
		if err := r.scheduleService.DeleteSchedule(req.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Jadwal dihapus"})
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
