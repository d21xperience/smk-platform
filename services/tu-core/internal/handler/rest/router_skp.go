package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerSKPRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kepegawaian/skp/templates", r.wrapWithContext(r.handleSKPTemplates))
	mux.HandleFunc("/api/v1/kepegawaian/skp/periodes", r.wrapWithContext(r.handleSKPPeriodeList))
	mux.HandleFunc("/api/v1/kepegawaian/skp/periodes/", r.wrapWithContext(r.handleSKPPeriodeByID))
	mux.HandleFunc("/api/v1/kepegawaian/skp/penilaians/", r.wrapWithContext(r.handleSKPPenilaianByID))
}

func (r *Router) handleSKPTemplates(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	templates, err := r.skpService.ListTemplates(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, templates)
}

func (r *Router) handleSKPPeriodeList(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		list, err := r.skpService.ListPeriodes(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		var body repository.SKPPeriode
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		id, err := r.skpService.CreatePeriode(req.Context(), body)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Periode berhasil dibuat", "id": id})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleSKPPeriodeByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/skp/periodes/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	periodeID := parts[0]

	if len(parts) == 2 && parts[1] == "penilaians" {
		if req.Method == http.MethodGet {
			list, err := r.skpService.ListPenilaianByPeriode(req.Context(), periodeID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, list)
		} else {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
		return
	}

	writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL format"})
}

func (r *Router) handleSKPPenilaianByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/skp/penilaians/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	id := parts[0]

	// GET Detail Penilaian
	if len(parts) == 1 && req.Method == http.MethodGet {
		penilaian, hasil, detail, err := r.skpService.GetPenilaianDetail(req.Context(), id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"penilaian": penilaian,
			"hasil":     hasil,
			"detail":    detail,
		})
		return
	}

	// POST / PUT Actions
	if len(parts) == 2 {
		action := parts[1]
		var err error
		switch action {
		case "submit":
			if req.Method != http.MethodPost {
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
				return
			}
			var body struct {
				CatatanPenilai string                     `json:"catatanPenilai"`
				Detail         []repository.SKPDetailItem `json:"detail"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
				return
			}
			hasil := repository.SKPHasil{CatatanPenilai: body.CatatanPenilai}
			err = r.skpService.SubmitPenilaian(req.Context(), id, hasil, body.Detail)

		case "approve":
			if req.Method != http.MethodPut {
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
				return
			}
			var body struct {
				Catatan string `json:"catatan"`
			}
			_ = json.NewDecoder(req.Body).Decode(&body)
			err = r.skpService.ApprovePenilaian(req.Context(), id, body.Catatan, opCtx.SchoolID)

		case "return":
			if req.Method != http.MethodPut {
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
				return
			}
			var body struct {
				Catatan string `json:"catatan"`
			}
			_ = json.NewDecoder(req.Body).Decode(&body)
			err = r.skpService.ReturnPenilaian(req.Context(), id, body.Catatan)

		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unknown action"})
			return
		}

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Aksi " + action + " berhasil"})
		return
	}

	writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL format"})
}
