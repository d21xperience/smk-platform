package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/service"
)

type ArchiveHandler struct {
	svc *service.ArchiveService
}

func NewArchiveHandler(svc *service.ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{svc: svc}
}

func (h *ArchiveHandler) HandleCollection(w http.ResponseWriter, req *http.Request) {
	// POST /api/v1/archives (Arsipkan Surat)
	// Body: { "correspondence_id": "...", "lokasi_rak": "...", ... }
	var body struct {
		CorrespondenceID string `json:"correspondence_id"`
		service.ArchiveCorrespondenceRequest
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	res, err := h.svc.ArsipkanSurat(req.Context(), body.CorrespondenceID, body.ArchiveCorrespondenceRequest)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (h *ArchiveHandler) HandleByID(w http.ResponseWriter, req *http.Request) {
	// Path: /api/v1/archives/{id}/{action}
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/archives/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if len(parts) == 1 {
		if req.Method == http.MethodGet {
			// Endpoint untuk get archive by ID (jika dibutuhkan)
			writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
			return
		}
	}

	action := parts[1]
	switch action {
	case "borrow":
		var body service.BorrowArchiveRequest
		_ = json.NewDecoder(req.Body).Decode(&body)
		if err := h.svc.PinjamArsip(req.Context(), id, body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Arsip berhasil dipinjam"})

	case "return":
		if err := h.svc.KembalikanArsip(req.Context(), id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Arsip berhasil dikembalikan"})

	case "move":
		var body service.MoveArchiveRequest
		_ = json.NewDecoder(req.Body).Decode(&body)
		if err := h.svc.PindahArsip(req.Context(), id, body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Arsip berhasil dipindahkan"})

	case "destroy":
		if err := h.svc.MusnahkanArsip(req.Context(), id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Arsip berhasil dimusnahkan"})

	case "movements":
		res, err := h.svc.ListMovements(req.Context(), id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, res)

	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Action not found"})
	}
}
