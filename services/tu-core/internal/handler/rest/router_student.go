package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerStudentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/students", r.handleStudents)
	mux.HandleFunc("/api/v1/students/available", r.handleAvailableStudents)
	mux.HandleFunc("/api/v1/students/", r.handleStudentByID)
	mux.HandleFunc("/api/v1/students/{id}/history", r.handleStudentHistory)
}

func (r *Router) handleStudents(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}

	switch req.Method {
	case http.MethodPost:
		r.createStudent(w, req.WithContext(ctx))
	case http.MethodGet:
		r.listStudents(w, req.WithContext(ctx))
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) createStudent(w http.ResponseWriter, req *http.Request) {
	var reqBody service.CreateStudentRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	student, err := r.studentService.CreateStudent(
		req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID,
		opCtx.Semester, opCtx.UserID, opCtx.Role, reqBody,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else if errors.Is(err, service.ErrDuplicateEntry) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusCreated, student)
}

func (r *Router) listStudents(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())

	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(req.URL.Query().Get("offset"))
	searchStr := req.URL.Query().Get("search")
	statusStr := req.URL.Query().Get("status")
	classID := req.URL.Query().Get("class_id")

	if limit <= 0 { limit = 20 }
	if offset < 0 { offset = 0 }

	filter := repository.StudentFilter{
		Limit: limit, Offset: offset, Keyword: searchStr, Status: statusStr, ClassID: classID,
	}

	students, err := r.studentService.ListStudentsForUI(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, students)
}

func (r *Router) handleStudentByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/students/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	studentID := parts[0]

	opCtx := middleware.GetOperationalContext(ctx)

	switch req.Method {
	case http.MethodGet:
		student, err := r.studentService.GetStudentForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, student)

	case http.MethodPut:
		// PERBAIKAN: Sesuaikan dengan payload JSON dari Frontend (camelCase)
		var reqBody struct {
			Name        string `json:"name"`
			NIS         string `json:"nis"`
			NISN        string `json:"nisn"`
			BirthPlace  string `json:"birthPlace"`
			BirthDate   string `json:"birthDate"`
			Gender      string `json:"gender"`
			Address     string `json:"address"`
			Phone       string `json:"phone"`
			ClassName   string `json:"classId"`
			ParentName  string `json:"parentName"`
			ParentPhone string `json:"parentPhone"`
		}
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body: " + err.Error()})
			return
		}
		student, err := r.studentService.UpdateStudentForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID, reqBody)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, student)

	case http.MethodDelete:
		err := r.studentService.DeleteStudentForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleAvailableStudents(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}

	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	opCtx := middleware.GetOperationalContext(ctx)
	keyword := req.URL.Query().Get("keyword")
	filter := repository.StudentFilter{Keyword: keyword, Limit: 100, Offset: 0}

	resp, err := r.studentService.ListAvailableStudentsForUI(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
