package rest

import (
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// registerKeuanganPiutangRoutes — endpoint piutang & aging.
//
// Route ordering: Go 1.22+ ServeMux memilih pattern paling spesifik,
// sehingga /aging dan /per-siswa otomatis diutamakan daripada /{nis}.
func (r *Router) registerKeuanganPiutangRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeTransaksiBasePath+"/piutang", r.handlePiutangCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/piutang/aging", r.handlePiutangAging)
	mux.HandleFunc(financeTransaksiBasePath+"/piutang/per-siswa", r.handlePiutangPerSiswa)
	mux.HandleFunc(financeTransaksiBasePath+"/piutang/siswa/{nis}", r.handlePiutangByNIS)
}

// =====================================================================
// Piutang Collection — GET list flat (per tagihan)
// =====================================================================
func (r *Router) handlePiutangCollection(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	periodID := req.URL.Query().Get("academic_period_id")
	if periodID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "academic_period_id wajib diisi"})
		return
	}
	filter := models.PiutangFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: periodID,
		NIS:              req.URL.Query().Get("nis"),
		KelasNama:        req.URL.Query().Get("kelas_nama"),
		AgingGroup:       req.URL.Query().Get("aging_group"),
		JenisIuranID:     req.URL.Query().Get("jenis_iuran_id"),
		Limit:            parseIntQuery(req, "limit", 50),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.piutangService.ListPiutang(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    items,
		"pagination": map[string]int{
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}

// =====================================================================
// Piutang Aging — GET aging report (summary buckets)
// =====================================================================
func (r *Router) handlePiutangAging(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	periodID := req.URL.Query().Get("academic_period_id")
	if periodID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "academic_period_id wajib diisi"})
		return
	}
	report, err := r.piutangService.GetAgingReport(req.Context(), schoolID, periodID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": report})
}

// =====================================================================
// Piutang Per-Siswa — GET grouped list (setiap siswa 1 row + items)
// =====================================================================
func (r *Router) handlePiutangPerSiswa(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	periodID := req.URL.Query().Get("academic_period_id")
	if periodID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "academic_period_id wajib diisi"})
		return
	}
	filter := models.PiutangFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: periodID,
		NIS:              req.URL.Query().Get("nis"),
		KelasNama:        req.URL.Query().Get("kelas_nama"),
		Limit:            parseIntQuery(req, "limit", 50),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.piutangService.ListPiutangPerSiswa(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    items,
		"pagination": map[string]int{
			"total":  total,
			"limit":  filter.Limit,
			"offset": filter.Offset,
		},
	})
}

// =====================================================================
// Piutang By NIS — GET detail 1 siswa (grouped)
// =====================================================================
func (r *Router) handlePiutangByNIS(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	periodID := req.URL.Query().Get("academic_period_id")
	if periodID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "academic_period_id wajib diisi"})
		return
	}
	nis := req.PathValue("nis")
	if nis == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nis wajib diisi"})
		return
	}
	item, err := r.piutangService.GetPiutangByNIS(req.Context(), schoolID, periodID, nis)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}
