package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/service"
)

type ExpeditionHandler struct {
	svc *service.ExpeditionService
}

func NewExpeditionHandler(svc *service.ExpeditionService) *ExpeditionHandler {
	return &ExpeditionHandler{svc: svc}
}

func (h *ExpeditionHandler) Create(w http.ResponseWriter, req *http.Request, correspondenceID string) {
	var body service.CreateExpeditionRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	res, err := h.svc.CreateExpedition(req.Context(), correspondenceID, body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *ExpeditionHandler) ListByCorrespondence(w http.ResponseWriter, req *http.Request, correspondenceID string) {
	res, err := h.svc.GetByCorrespondenceID(req.Context(), correspondenceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *ExpeditionHandler) HandleCollection(w http.ResponseWriter, req *http.Request) {
	// Endpoint /api/v1/correspondences/{id}/expeditions
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/")
	id := strings.Split(path, "/")[0]

	switch req.Method {
	case http.MethodPost:
		h.Create(w, req, id)
	case http.MethodGet:
		h.ListByCorrespondence(w, req, id)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
