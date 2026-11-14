package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/service"
)

const financeTransaksiBasePath = "/api/v1/keuangan"

// registerKeuanganTransaksiRoutes — mendaftarkan endpoint tagihan,
// penerimaan, dan BKU.
func (r *Router) registerKeuanganTransaksiRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeTransaksiBasePath+"/tagihan", r.handleTagihanCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/tagihan/{id}", r.handleTagihanItem)
	mux.HandleFunc(financeTransaksiBasePath+"/penerimaan", r.handlePenerimaanCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/penerimaan/{no_bkm...}", r.handlePenerimaanItem)
	mux.HandleFunc(financeTransaksiBasePath+"/bku", r.handleBKUCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/bku/saldo", r.handleBKUSaldo)
}

// =====================================================================
// TAGIHAN
// =====================================================================
func (r *Router) handleTagihanCollection(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := models.TagihanFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		NIS:              req.URL.Query().Get("nis"),
		JenisIuranID:     req.URL.Query().Get("jenis_iuran_id"),
		Status:           req.URL.Query().Get("status"),
		Bulan:            parseIntQuery(req, "bulan", 0),
		Tahun:            parseIntQuery(req, "tahun", 0),
		Limit:            parseIntQuery(req, "limit", 50),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.iuranService.ListTagihan(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

func (r *Router) handleTagihanItem(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	id := req.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id wajib diisi"})
		return
	}
	item, err := r.iuranService.GetTagihan(req.Context(), schoolID, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

// =====================================================================
// PENERIMAAN (auto BKM + auto BKU + auto update tagihan)
// =====================================================================
func (r *Router) handlePenerimaanCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listPenerimaan(w, req, schoolID)
	case http.MethodPost:
		r.createPenerimaan(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listPenerimaan(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.PenerimaanFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		NIS:              req.URL.Query().Get("nis"),
		JenisIuranID:     req.URL.Query().Get("jenis_iuran_id"),
		TanggalDari:      req.URL.Query().Get("tanggal_dari"),
		TanggalSampai:    req.URL.Query().Get("tanggal_sampai"),
		Limit:            parseIntQuery(req, "limit", 50),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.iuranService.ListPenerimaan(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

func (r *Router) createPenerimaan(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload service.TerimaPembayaranRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if payload.AcademicPeriodID == "" {
		payload.AcademicPeriodID = req.URL.Query().Get("academic_period_id")
	}
	resp, err := r.iuranService.TerimaPembayaran(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": resp})
}

func (r *Router) handlePenerimaanItem(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	noBKM := req.PathValue("no_bkm")
	if noBKM == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no_bkm wajib diisi"})
		return
	}
	item, err := r.iuranService.GetPenerimaan(req.Context(), schoolID, noBKM)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

// =====================================================================
// BUKU KAS UMUM
// =====================================================================
func (r *Router) handleBKUCollection(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := models.BukuKasFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		KodeRekeningID:   req.URL.Query().Get("kode_rekening_id"),
		TanggalDari:      req.URL.Query().Get("tanggal_dari"),
		TanggalSampai:    req.URL.Query().Get("tanggal_sampai"),
		Limit:            parseIntQuery(req, "limit", 100),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.bukuKasService.ListBKU(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

func (r *Router) handleBKUSaldo(w http.ResponseWriter, req *http.Request) {
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
	tanggal := req.URL.Query().Get("tanggal")
	saldo, tgl, err := r.bukuKasService.GetSaldo(req.Context(), schoolID, periodID, tanggal)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"saldo":   saldo,
			"tanggal": tgl,
		},
	})
}
