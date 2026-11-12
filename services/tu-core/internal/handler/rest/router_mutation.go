package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerMutationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/mutations", r.handleMutations)
	mux.HandleFunc("/api/v1/mutations/", r.handleMutationByID)
}

func (r *Router) handleMutations(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodPost:
		r.createMutation(w, req)
	case http.MethodGet:
		r.listMutations(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) createMutation(w http.ResponseWriter, req *http.Request) {
	var reqBody struct {
		StudentID     *string `json:"studentId"`
		MutationType  string  `json:"mutationType"`
		NIS           string  `json:"nis"`
		NISN          string  `json:"nisn"`
		Nama          string  `json:"nama"`
		JenisKelamin  string  `json:"jenisKelamin"`
		TempatLahir   string  `json:"tempatLahir"`
		TanggalLahir  *string `json:"tanggalLahir"`
		AsalSekolah   string  `json:"asalSekolah"`
		TujuanSekolah string  `json:"tujuanSekolah"`
		Alasan        string  `json:"alasan"`
		DokumenURL    string  `json:"dokumenUrl"`
	}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	var tglLahir *time.Time
	if reqBody.TanggalLahir != nil && *reqBody.TanggalLahir != "" {
		if t, err := time.Parse("2006-01-02", *reqBody.TanggalLahir); err == nil {
			tglLahir = &t
		}
	}

	serviceReq := service.CreateMutationRequest{
		StudentID:     reqBody.StudentID,
		MutationType:  reqBody.MutationType,
		NIS:           reqBody.NIS,
		NISN:          reqBody.NISN,
		Nama:          reqBody.Nama,
		JenisKelamin:  reqBody.JenisKelamin,
		TempatLahir:   reqBody.TempatLahir,
		TanggalLahir:  tglLahir,
		AsalSekolah:   reqBody.AsalSekolah,
		TujuanSekolah: reqBody.TujuanSekolah,
		Alasan:        reqBody.Alasan,
		DokumenURL:    reqBody.DokumenURL,
	}

	m, err := r.mutationService.CreateMutation(req.Context(), serviceReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, m)
}

func (r *Router) listMutations(w http.ResponseWriter, req *http.Request) {
	status := req.URL.Query().Get("status")
	mType := req.URL.Query().Get("type")

	mutations, err := r.mutationService.ListMutations(req.Context(), status, mType)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  mutations,
		"total": len(mutations),
	})
}

func (r *Router) handleMutationByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/mutations/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}

	mutationID := parts[0]

	// Sub-path: /{id}/approve atau /{id}/reject
	if len(parts) > 1 {
		action := parts[1]
		switch action {
		case "approve":
			if req.Method != http.MethodPost {
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
				return
			}
			var body struct {
				Catatan string `json:"catatan"`
			}
			_ = json.NewDecoder(req.Body).Decode(&body)
			if err := r.mutationService.ApproveMutation(req.Context(), mutationID, body.Catatan); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Mutasi disetujui"})
			return

		case "reject":
			if req.Method != http.MethodPost {
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
				return
			}
			var body struct {
				Catatan string `json:"catatan"`
			}
			_ = json.NewDecoder(req.Body).Decode(&body)
			if err := r.mutationService.RejectMutation(req.Context(), mutationID, body.Catatan); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Mutasi ditolak"})
			return
		}
	}

	// GET /{id}
	if req.Method == http.MethodGet {
		m, err := r.mutationService.GetMutationByID(req.Context(), mutationID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, m)
		return
	}

	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
}

var _ = middleware.GetOperationalContext
