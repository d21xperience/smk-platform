package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

type ContextHandler struct {
	prefSvc *service.PreferenceService
}

func NewContextHandler(prefSvc *service.PreferenceService) *ContextHandler {
	return &ContextHandler{prefSvc: prefSvc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func (h *ContextHandler) GetPeriods(w http.ResponseWriter, req *http.Request) {
	periods, err := h.prefSvc.ListPeriods(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, periods)
}

func (h *ContextHandler) CreatePeriod(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Konteks operasional tidak ditemukan."})
		return
	}

	role := strings.ToLower(strings.TrimSpace(opCtx.Role))
	if role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak. Role Anda: " + role})
		return
	}

	var body struct {
		Year string `json:"year"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	if err := h.prefSvc.CreateAcademicYear(req.Context(), body.Year); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Tahun ajaran dan 2 semester berhasil ditambahkan"})
}

func (h *ContextHandler) ActivatePeriod(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Konteks operasional tidak ditemukan."})
		return
	}

	role := strings.ToLower(strings.TrimSpace(opCtx.Role))
	if role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak."})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/academic/periods/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "activate" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Path tidak valid"})
		return
	}

	if err := h.prefSvc.ActivatePeriod(req.Context(), parts[0]); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Periode berhasil diaktifkan"})
}

func (h *ContextHandler) DeletePeriod(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Konteks operasional tidak ditemukan."})
		return
	}

	role := strings.ToLower(strings.TrimSpace(opCtx.Role))
	if role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak."})
		return
	}

	// Path: /api/v1/academic/periods/{id}
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/academic/periods/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID periode tidak valid"})
		return
	}

	if err := h.prefSvc.DeletePeriod(req.Context(), parts[0]); err != nil {
		// Menangkap error jika semester aktif (Status 400 Bad Request)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Periode berhasil dihapus"})
}

func (h *ContextHandler) GetPreference(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("x-user-id")
	schoolID := req.Header.Get("x-school-id")
	semesterID, err := h.prefSvc.GetUserPreference(req.Context(), userID, schoolID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"semester_id": semesterID})
}

func (h *ContextHandler) UpdatePreference(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("x-user-id")
	schoolID := req.Header.Get("x-school-id")
	var body struct {
		SemesterID string `json:"semester_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if err := h.prefSvc.SaveUserPreference(req.Context(), userID, schoolID, body.SemesterID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *ContextHandler) TogglePSB(w http.ResponseWriter, req *http.Request) {
	var body struct {
		SemesterID string `json:"semester_id"`
		IsActive   bool   `json:"is_active"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}

	if err := h.prefSvc.TogglePSB(req.Context(), body.SemesterID, body.IsActive); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Status PSB berhasil diubah"})
}
