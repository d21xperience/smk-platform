package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerScholarshipRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/scholarships", r.handleScholarships)
	mux.HandleFunc("/api/v1/scholarships/", r.handleScholarshipByID)
}

func (r *Router) handleScholarships(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodGet:
		r.listScholarships(w, req)
	case http.MethodPost:
		r.createScholarship(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) createScholarship(w http.ResponseWriter, req *http.Request) {
	var reqBody service.CreateScholarshipRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body: " + err.Error()})
		return
	}

	result, err := r.scholarshipService.CreateScholarship(req.Context(), reqBody)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (r *Router) listScholarships(w http.ResponseWriter, req *http.Request) {
	status := req.URL.Query().Get("status")
	search := req.URL.Query().Get("search")

	items, err := r.scholarshipService.ListScholarships(req.Context(), status, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  items,
		"total": len(items),
	})
}

func (r *Router) handleScholarshipByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/scholarships/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	scholarshipID := parts[0]

	// Sub-path: /{id}/toggle-document
	if len(parts) > 1 && parts[1] == "toggle-document" {
		if req.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}
		var body struct {
			IsValid bool `json:"isValid"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.scholarshipService.ToggleDocumentValid(req.Context(), scholarshipID, body.IsValid); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true})
		return
	}

	switch req.Method {
	case http.MethodGet:
		item, err := r.scholarshipService.GetScholarshipByID(req.Context(), scholarshipID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, item)

	case http.MethodPut:
		var reqBody service.UpdateScholarshipRequest
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		result, err := r.scholarshipService.UpdateScholarship(req.Context(), scholarshipID, reqBody)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, result)

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
