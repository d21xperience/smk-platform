package rest

import (
	"encoding/json"
	"net/http"
	"sekolah-platform/services/tu-core/internal/middleware"
	"strings"
)

func (r *Router) registerPelayananRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kepegawaian/pelayanan/katalog", r.wrapWithContext(r.handleKatalogLayanan))
	mux.HandleFunc("/api/v1/kepegawaian/pelayanan/pengajuan", r.wrapWithContext(r.handlePengajuanList))
	mux.HandleFunc("/api/v1/kepegawaian/pelayanan/pengajuan/", r.wrapWithContext(r.handlePengajuanByID))
}

func (r *Router) handleKatalogLayanan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	list, err := r.pelayananService.ListKatalog(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handlePengajuanList(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	switch req.Method {
	case http.MethodGet:
		status := req.URL.Query().Get("status")
		list, err := r.pelayananService.ListPengajuan(req.Context(), status)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		// Endpoint ini dipanggil oleh VPS/siakad-guru atau manual
		var body struct {
			PtkID       string `json:"ptkId"`
			LayananKode string `json:"layananKode"`
			Tujuan      string `json:"tujuan"`
			Kepada      string `json:"kepada"`
			Keterangan  string `json:"keterangan"`
			DokURL      string `json:"dokumenUrl"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}

		id, err := r.pelayananService.CreatePengajuan(req.Context(), body.PtkID, body.LayananKode, body.Tujuan, body.Kepada, body.Keterangan, body.DokURL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Pengajuan berhasil", "id": id})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handlePengajuanByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/pelayanan/pengajuan/")
	parts := splitPath(path)
	if len(parts) < 2 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL"})
		return
	}
	id := parts[0]
	action := parts[1]

	var err error
	switch action {
	case "verify":
		var body struct{ Action, Catatan string }
		_ = json.NewDecoder(req.Body).Decode(&body)
		err = r.pelayananService.Verify(req.Context(), id, body.Action, body.Catatan, opCtx.UserID)
	case "approve":
		_, err = r.pelayananService.Approve(req.Context(), id, opCtx.UserID, opCtx.SchoolID)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unknown action"})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Aksi berhasil"})
}
