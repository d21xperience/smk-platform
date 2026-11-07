package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerAcademicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/subjects", r.handleSubjects)
	mux.HandleFunc("/api/v1/subjects/", r.handleSubjectByID)

	mux.HandleFunc("/api/v1/schedules/options", r.wrapWithContext(r.scheduleHandler.GetOptions))
	mux.HandleFunc("/api/v1/schedules", r.wrapWithContext(r.handleSchedules))
	mux.HandleFunc("/api/v1/schedules/{id}", r.wrapWithContext(r.handleScheduleByID))

	mux.HandleFunc("/api/v1/kurikulum/statistics", r.handleKurikulumStats)
	mux.HandleFunc("/api/v1/pembelajaran", r.handlePembelajaran)
	mux.HandleFunc("/api/v1/statistics", r.handleStatistics)

	mux.HandleFunc("/api/v1/classes", r.handleClasses)
	mux.HandleFunc("/api/v1/classes/{id}", r.handleClassByID)
	mux.HandleFunc("/api/v1/classes/{id}/members", r.handleClassMembers)
}

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

func (r *Router) handleScheduleByID(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodDelete {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]
		if err := r.scheduleService.DeleteSchedule(req.Context(), id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Jadwal dihapus"})
		return
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
}

func (r *Router) handleStatistics(w http.ResponseWriter, req *http.Request) {
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
	stats, err := r.statisticsService.GetDashboardStats(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (r *Router) handleSubjects(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		subjects, err := r.kurikulumService.ListSubjects(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": subjects})
	case http.MethodPost:
		var body repository.Subject
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.kurikulumService.CreateSubject(req.Context(), &body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, body)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleSubjectByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/subjects/")
	pathParts := splitPath(path)
	if len(pathParts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}

	id, err := strconv.Atoi(pathParts[0])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid subject ID format"})
		return
	}

	switch req.Method {
	case http.MethodPut:
		var body repository.Subject
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		body.ID = id
		if err := r.kurikulumService.UpdateSubject(ctx, &body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, body)

	case http.MethodDelete:
		if err := r.kurikulumService.DeleteSubject(ctx, id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Mata pelajaran berhasil dihapus"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handlePembelajaran(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	opCtx := middleware.GetOperationalContext(ctx)

	switch req.Method {
	case http.MethodGet:
		classID := req.URL.Query().Get("classId")
		list, err := r.kurikulumService.ListPembelajaran(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
	case http.MethodPost:
		var body struct {
			ClassID   string `json:"classId"`
			SubjectID string `json:"subjectId"`
			TeacherID string `json:"teacherId"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.kurikulumService.CreatePembelajaran(ctx, body.ClassID, body.SubjectID, body.TeacherID, opCtx.AcademicPeriodID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Pembelajaran ditambahkan"})
	}
}

func (r *Router) handleKurikulumStats(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	opCtx := middleware.GetOperationalContext(ctx)

	stats, err := r.kurikulumService.GetDashboardStats(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, stats)
}
