package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/service"
)

// registerKeuanganPengeluaranRoutes — endpoint pengajuan & pengeluaran.
func (r *Router) registerKeuanganPengeluaranRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeTransaksiBasePath+"/pengajuan", r.handlePengajuanCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/pengajuan/{id}", r.handlePengajuanItem)
	mux.HandleFunc(financeTransaksiBasePath+"/pengajuan/{id}/verifikasi", r.handlePengajuanVerifikasi)
	mux.HandleFunc(financeTransaksiBasePath+"/pengajuan/{id}/approval", r.handlePengajuanApproval)
	mux.HandleFunc(financeTransaksiBasePath+"/pengeluaran", r.handlePengeluaranCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/pengeluaran/{no_bkk...}", r.handlePengeluaranItem)
}

// =====================================================================
// PENGAJUAN PENGELUARAN
// =====================================================================
func (r *Router) handlePengajuanCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listPengajuan(w, req, schoolID)
	case http.MethodPost:
		r.createPengajuan(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listPengajuan(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.PengajuanFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		Status:           req.URL.Query().Get("status"),
		PemohonUserID:    req.URL.Query().Get("pemohon_user_id"),
		Limit:            parseIntQuery(req, "limit", 50),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.pengeluaranService.ListPengajuan(req.Context(), filter)
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

func (r *Router) createPengajuan(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload service.AjukanPengeluaranRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if payload.AcademicPeriodID == "" {
		payload.AcademicPeriodID = req.URL.Query().Get("academic_period_id")
	}
	item, err := r.pengeluaranService.AjukanPengeluaran(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": item})
}

func (r *Router) handlePengajuanItem(w http.ResponseWriter, req *http.Request) {
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
	item, err := r.pengeluaranService.GetPengajuan(req.Context(), schoolID, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

func (r *Router) handlePengajuanVerifikasi(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
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
	var payload service.VerifikasiBendaharaRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	payload.PengajuanID = id
	item, err := r.pengeluaranService.VerifikasiBendahara(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

func (r *Router) handlePengajuanApproval(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
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
	var payload service.ApprovalKepalaRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	payload.PengajuanID = id
	resp, err := r.pengeluaranService.ApprovalKepala(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": resp})
}

// =====================================================================
// PENGELUARAN
// =====================================================================
func (r *Router) handlePengeluaranCollection(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	filter := models.PengeluaranFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		Status:           req.URL.Query().Get("status"),
		TanggalDari:      req.URL.Query().Get("tanggal_dari"),
		TanggalSampai:    req.URL.Query().Get("tanggal_sampai"),
		Limit:            parseIntQuery(req, "limit", 50),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.pengeluaranService.ListPengeluaran(req.Context(), filter)
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

func (r *Router) handlePengeluaranItem(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	noBKK := req.PathValue("no_bkk")
	if noBKK == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no_bkk wajib diisi"})
		return
	}
	item, err := r.pengeluaranService.GetPengeluaran(req.Context(), schoolID, noBKK)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}
