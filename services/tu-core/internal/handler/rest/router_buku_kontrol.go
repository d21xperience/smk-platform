package rest

import (
	"net/http"
	"strings"
)

func (r *Router) registerBukuKontrolRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kepegawaian/buku-kontrol/dashboard", r.wrapWithContext(r.handleBukuKontrolDashboard))
	mux.HandleFunc("/api/v1/kepegawaian/buku-kontrol/initialize", r.wrapWithContext(r.handleInitializeBukuKontrol))
	mux.HandleFunc("/api/v1/kepegawaian/buku-kontrol/golongan", r.wrapWithContext(r.handleListGolongan))
	mux.HandleFunc("/api/v1/kepegawaian/buku-kontrol/", r.wrapWithContext(r.handleBukuKontrolByPtkID))
}

func (r *Router) handleBukuKontrolDashboard(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	data, err := r.bukuKontrolService.GetDashboard(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (r *Router) handleInitializeBukuKontrol(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	count, err := r.bukuKontrolService.InitializeMassal(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"message": "Inisialisasi berhasil", "total": count})
}

func (r *Router) handleListGolongan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	list, err := r.bukuKontrolService.ListGolongan(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handleBukuKontrolByPtkID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	ptkID := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/buku-kontrol/")
	if ptkID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "PTK ID required"})
		return
	}
	data, err := r.bukuKontrolService.GetByPtkID(req.Context(), ptkID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Data tidak ditemukan"})
		return
	}
	writeJSON(w, http.StatusOK, data)
}
