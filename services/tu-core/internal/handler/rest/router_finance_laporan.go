package rest

import (
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// registerKeuanganLaporanRoutes — endpoint laporan keuangan (FASE 6).
//
// Semua endpoint butuh: school_id (query/header) + academic_period_id.
func (r *Router) registerKeuanganLaporanRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeTransaksiBasePath+"/laporan/neraca", r.handleLaporanNeraca)
	mux.HandleFunc(financeTransaksiBasePath+"/laporan/posisi-keuangan", r.handleLaporanPosisiKeuangan)
	mux.HandleFunc(financeTransaksiBasePath+"/laporan/lra", r.handleLaporanLRA)
	mux.HandleFunc(financeTransaksiBasePath+"/laporan/dashboard", r.handleLaporanDashboard)
	mux.HandleFunc(financeTransaksiBasePath+"/laporan/ringkasan-bulan", r.handleLaporanRingkasanBulan)
}

// =====================================================================
// Shared: extract LaporanFilter dari query string
// =====================================================================
func buildLaporanFilter(req *http.Request, schoolID string) models.LaporanFilter {
	return models.LaporanFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		Tanggal:          req.URL.Query().Get("tanggal"),
		Tahun:            parseIntQuery(req, "tahun", 0),
		SampaiBulan:      parseIntQuery(req, "sampai_bulan", 0),
		Bulan:            parseIntQuery(req, "bulan", 0),
		Mode:             req.URL.Query().Get("mode"),
		PenutupanID:      req.URL.Query().Get("penutupan_id"),
	}
}

// =====================================================================
// GET /laporan/neraca
// =====================================================================
func (r *Router) handleLaporanNeraca(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := buildLaporanFilter(req, schoolID)
	report, err := r.laporanService.GetNeraca(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": report})
}

// =====================================================================
// GET /laporan/posisi-keuangan (strict — harus balance)
// =====================================================================
func (r *Router) handleLaporanPosisiKeuangan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := buildLaporanFilter(req, schoolID)
	report, err := r.laporanService.GetPosisiKeuangan(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": report})
}

// =====================================================================
// GET /laporan/lra
// =====================================================================
func (r *Router) handleLaporanLRA(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := buildLaporanFilter(req, schoolID)
	report, err := r.laporanService.GetLra(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": report})
}

// =====================================================================
// GET /laporan/dashboard
// =====================================================================
func (r *Router) handleLaporanDashboard(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := buildLaporanFilter(req, schoolID)
	summary, err := r.laporanService.GetDashboard(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": summary})
}

// =====================================================================
// GET /laporan/ringkasan-bulan
// =====================================================================
func (r *Router) handleLaporanRingkasanBulan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := buildLaporanFilter(req, schoolID)
	summary, err := r.laporanService.GetRingkasanBulan(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": summary})
}
