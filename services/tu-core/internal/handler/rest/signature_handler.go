package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/service"
)

type SignatureHandler struct {
	svc *service.SignatureService
}

func NewSignatureHandler(svc *service.SignatureService) *SignatureHandler {
	return &SignatureHandler{svc: svc}
}

// SignCorrespondence (POST /api/v1/correspondences/{id}/sign)
func (h *SignatureHandler) Sign(w http.ResponseWriter, req *http.Request, correspondenceID string) {
	var body service.SignCorrespondenceRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	res, err := h.svc.SignCorrespondence(req.Context(), correspondenceID, body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// VerifySignature (GET /api/v1/signatures/verify/{code})
func (h *SignatureHandler) Verify(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/signatures/verify/")
	if path == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing verification code"})
		return
	}

	res, err := h.svc.VerifyCode(req.Context(), path)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}
