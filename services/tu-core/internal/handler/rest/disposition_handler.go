package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

type DispositionHandler struct {
	svc *service.DispositionService
}

func NewDispositionHandler(svc *service.DispositionService) *DispositionHandler {
	return &DispositionHandler{svc: svc}
}

// CreateDisposition (POST /api/v1/correspondences/{id}/dispositions)
func (h *DispositionHandler) Create(w http.ResponseWriter, req *http.Request, correspondenceID string) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	var body struct {
		ToUserID string `json:"to_user_id"`
		Note     string `json:"note"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	svcReq := service.CreateDispositionRequest{
		CorrespondenceID: correspondenceID,
		ToUserID:         body.ToUserID,
		Note:             body.Note,
	}

	res, err := h.svc.CreateDisposition(req.Context(), svcReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

// ListByCorrespondence (GET /api/v1/correspondences/{id}/dispositions)
func (h *DispositionHandler) ListByCorrespondence(w http.ResponseWriter, req *http.Request, correspondenceID string) {
	res, err := h.svc.ListByCorrespondenceID(req.Context(), correspondenceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// GetMyInbox (GET /api/v1/dispositions/inbox)
func (h *DispositionHandler) GetMyInbox(w http.ResponseWriter, req *http.Request) {
	res, err := h.svc.GetMyInbox(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// UpdateStatus (PUT /api/v1/dispositions/{id}/status)
func (h *DispositionHandler) UpdateStatus(w http.ResponseWriter, req *http.Request, dispositionID string) {
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	svcReq := service.UpdateDispositionStatusRequest{Status: body.Status}
	if err := h.svc.UpdateStatus(req.Context(), dispositionID, svcReq); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Status disposisi diperbarui"})
}

// Helper untuk ekstrak ID dari path
func extractIDFromPath(path, prefix string) string {
	p := strings.TrimPrefix(path, prefix)
	parts := strings.Split(p, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
