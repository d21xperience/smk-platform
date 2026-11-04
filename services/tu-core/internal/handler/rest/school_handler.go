package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

type SchoolHandler struct {
	schoolSvc *service.SchoolService
}

func NewSchoolHandler(schoolSvc *service.SchoolService) *SchoolHandler {
	return &SchoolHandler{schoolSvc: schoolSvc}
}

func (h *SchoolHandler) GetSchool(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	school, err := h.schoolSvc.GetSchool(req.Context(), opCtx.SchoolID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, school)
}

func (h *SchoolHandler) UpdateSchool(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	role := strings.ToLower(strings.TrimSpace(opCtx.Role))
	if role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak. Hanya Admin."})
		return
	}

	var body repository.School
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	body.SchoolID = opCtx.SchoolID // Paksa ID sesuai konteks operasional
	if err := h.schoolSvc.UpdateSchool(req.Context(), &body); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Profil sekolah berhasil diperbarui"})
}
