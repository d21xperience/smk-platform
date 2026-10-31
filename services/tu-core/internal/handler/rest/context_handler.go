package rest

import (
	"encoding/json"
	"net/http"

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
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	periods, err := h.prefSvc.ListPeriods(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, periods)
}

func (h *ContextHandler) GetPreference(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("x-user-id")
	schoolID := req.Header.Get("x-school-id")
	if userID == "" || schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "x-user-id and x-school-id headers are required"})
		return
	}
	semesterID, err := h.prefSvc.GetUserPreference(req.Context(), userID, schoolID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"user_id":     userID,
		"school_id":   schoolID,
		"semester_id": semesterID,
	})
}

func (h *ContextHandler) UpdatePreference(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("x-user-id")
	schoolID := req.Header.Get("x-school-id")
	if userID == "" || schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "x-user-id and x-school-id headers are required"})
		return
	}
	var body struct {
		SemesterID string `json:"semester_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if body.SemesterID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "semester_id is required"})
		return
	}
	if err := h.prefSvc.SaveUserPreference(req.Context(), userID, schoolID, body.SemesterID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}
