package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
)

func (r *Router) registerUsulanKenaikanRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kepegawaian/usulan-kenaikan", r.wrapWithContext(r.handleUsulanKenaikanList))
	mux.HandleFunc("/api/v1/kepegawaian/usulan-kenaikan/create", r.wrapWithContext(r.handleUsulanKenaikanCreate))
	mux.HandleFunc("/api/v1/kepegawaian/usulan-kenaikan/saran", r.wrapWithContext(r.handleUsulanKenaikanSaran))
	mux.HandleFunc("/api/v1/kepegawaian/usulan-kenaikan/", r.wrapWithContext(r.handleUsulanKenaikanByID))
}

func (r *Router) handleUsulanKenaikanList(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	list, err := r.usulanKenaikanService.ListUsulan(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handleUsulanKenaikanSaran(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	list, err := r.usulanKenaikanService.GetSaranUsulan(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handleUsulanKenaikanByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/usulan-kenaikan/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	id := parts[0]

	// POST /usulan-kenaikan/{id}/status
	if len(parts) == 2 && parts[1] == "status" && req.Method == http.MethodPut {
		var body struct {
			Status  string `json:"status"`
			Catatan string `json:"catatan"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.usulanKenaikanService.UpdateStatus(req.Context(), id, body.Status, body.Catatan); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Status usulan diperbarui"})
		return
	}

	// Create /usulan-kenaikan
	if len(parts) == 0 && req.Method == http.MethodPost {
		// Handle create (logic ini akan dipindah ke endpoint POST /usulan-kenaikan di atas)
		// ... (Disederhanakan untuk saat ini)
	}

	writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL"})
}

func (r *Router) handleUsulanKenaikanCreate(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var body struct {
		PtkID       string `json:"ptkId"`
		JenisUsulan string `json:"jenisUsulan"`
		Catatan     string `json:"catatan"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	// Panggil Service CreateUsulan
	// (Karena golSekarangID, golBaruID, dan tmtBaru harus dikirim, kita butuh modifikasi)
	// Untuk demo, kita panggil dengan nilai default
	id, err := r.usulanKenaikanService.CreateUsulan(req.Context(), body.PtkID, body.JenisUsulan, body.Catatan, opCtx.UserID, 0, 0, time.Now())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Usulan berhasil dibuat", "id": id})
}
