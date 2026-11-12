package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
)

func (r *Router) handleJurusan(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	switch req.Method {
	case http.MethodGet:
		list, err := r.jurusanService.List(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})

	case http.MethodPost:
		role := strings.ToLower(strings.TrimSpace(opCtx.Role))
		if role != "admin" {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak."})
			return
		}
		var body struct {
			JurusanID   string `json:"jurusan_id"`
			NamaJurusan string `json:"nama_jurusan"`
			LevelBidang string `json:"level_bidang_id"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
			return
		}
		if err := r.jurusanService.Create(req.Context(), body.JurusanID, body.NamaJurusan, body.LevelBidang); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Jurusan berhasil ditambahkan"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
