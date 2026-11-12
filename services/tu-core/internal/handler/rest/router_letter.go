package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerLetterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/letters", r.handleLetters)
	mux.HandleFunc("/api/v1/letters/", r.handleLetterByID)
}

func (r *Router) handleLetters(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodGet:
		r.listLetters(w, req)
	case http.MethodPost:
		r.createLetter(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) createLetter(w http.ResponseWriter, req *http.Request) {
	var reqBody service.CreateLetterRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body: " + err.Error()})
		return
	}

	letter, err := r.letterService.CreateLetter(req.Context(), reqBody)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, letter)
}

func (r *Router) listLetters(w http.ResponseWriter, req *http.Request) {
	status := req.URL.Query().Get("status")
	search := req.URL.Query().Get("search")

	letters, err := r.letterService.ListLetters(req.Context(), status, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  letters,
		"total": len(letters),
	})
}

func (r *Router) handleLetterByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/letters/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	letterID := parts[0]

	// Sub-path: /{id}/complete
	if len(parts) > 1 && parts[1] == "complete" {
		if req.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}
		if err := r.letterService.CompleteLetter(req.Context(), letterID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Surat ditandai selesai"})
		return
	}

	switch req.Method {
	case http.MethodGet:
		letter, err := r.letterService.GetLetterByID(req.Context(), letterID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, letter)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
