package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerPSBRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/registrations", r.handleRegistrations)
	mux.HandleFunc("/api/v1/registrations/generate-nis", r.handleGenerateNISMassal)
	mux.HandleFunc("/api/v1/registrations/", r.handleRegistrationByID)
}

func (r *Router) handleRegistrations(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodPost:
		r.createRegistration(w, req)
	case http.MethodGet:
		r.listRegistrations(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) createRegistration(w http.ResponseWriter, req *http.Request) {
	var reqBody service.CreateRegistrationRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
		return
	}

	reg, err := r.registrationService.CreateRegistration(req.Context(), reqBody)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, reg)
}

func (r *Router) listRegistrations(w http.ResponseWriter, req *http.Request) {
	status := req.URL.Query().Get("status")
	search := req.URL.Query().Get("search")

	regs, err := r.registrationService.ListRegistrations(req.Context(), status, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  regs,
		"total": len(regs),
	})
}

func (r *Router) handleGenerateNISMassal(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}

	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var reqBody struct {
		Prefix string `json:"prefix"`
	}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	// PERBAIKAN: Gunakan `ctx` (hasil enrichment), bukan req.Context()
	count, err := r.registrationService.GenerateNISMassal(ctx, reqBody.Prefix)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "NIS berhasil di-generate",
		"count":   count,
	})
}

func (r *Router) handleRegistrationByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/registrations/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}

	regID := parts[0]

	// Cek sub-path: /{id}/activate
	if len(parts) > 1 && parts[1] == "activate" {
		if req.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}
		// PERBAIKAN: Gunakan `ctx` (hasil enrichment)
		if err := r.registrationService.ActivateRegistration(ctx, regID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Siswa berhasil diaktifkan"})
		return
	}

	// GET /{id}
	switch req.Method {
	case http.MethodGet:
		// PERBAIKAN: Gunakan `ctx` (hasil enrichment)
		reg, err := r.registrationService.GetRegistrationByID(ctx, regID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, reg)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
