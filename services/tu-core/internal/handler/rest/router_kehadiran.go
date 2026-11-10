package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func (r *Router) registerKehadiranRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kepegawaian/kehadiran/harian", r.wrapWithContext(r.handleKehadiranHarian))
	mux.HandleFunc("/api/v1/kepegawaian/kehadiran/manual", r.wrapWithContext(r.handleKehadiranManual))
	mux.HandleFunc("/api/v1/kepegawaian/kehadiran/collect", r.wrapWithContext(r.handleKehadiranCollect))
	mux.HandleFunc("/api/v1/kepegawaian/kehadiran/rekap-bulanan", r.wrapWithContext(r.handleRekapBulanan))
	mux.HandleFunc("/api/v1/kepegawaian/kehadiran/rekap-bulanan/generate", r.wrapWithContext(r.handleGenerateRekap))
}

func (r *Router) handleKehadiranHarian(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	tglStr := req.URL.Query().Get("tanggal")
	tgl := time.Now()
	if tglStr != "" {
		if t, err := time.Parse("2006-01-02", tglStr); err == nil {
			tgl = t
		}
	}

	list, err := r.kehadiranService.GetAttendanceByDate(req.Context(), tgl)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handleKehadiranManual(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var body struct {
		PtkID      string  `json:"ptkId"`
		Tanggal    string  `json:"tanggal"`
		Status     string  `json:"status"`
		JamMasuk   *string `json:"jamMasuk"`
		JamPulang  *string `json:"jamPulang"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	tgl, err := time.Parse("2006-01-02", body.Tanggal)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format tanggal harus YYYY-MM-DD"})
		return
	}

	if err := r.kehadiranService.ManualInput(req.Context(), body.PtkID, tgl, body.Status, body.JamMasuk, body.JamPulang); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Kehadiran berhasil disimpan"})
}

func (r *Router) handleKehadiranCollect(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	tglStr := req.URL.Query().Get("tanggal")
	tgl := time.Now()
	if tglStr != "" {
		if t, err := time.Parse("2006-01-02", tglStr); err == nil {
			tgl = t
		}
	}

	count, err := r.kehadiranService.AutoCollect(req.Context(), tgl)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"message": "Auto-collect selesai", "total": count})
}

func (r *Router) handleRekapBulanan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	list, err := r.kehadiranService.GetRekapBulanan(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handleGenerateRekap(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var body struct {
		Bulan int `json:"bulan"`
		Tahun int `json:"tahun"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	id, err := r.kehadiranService.GenerateRekapBulanan(req.Context(), body.Bulan, body.Tahun)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "Rekap berhasil dibuat", "id": id})
}

// Helper untuk string pointer
func strPtr(s string) *string {
	if s == "" { return nil }
	return &s
}

var _ = strings.TrimSpace
