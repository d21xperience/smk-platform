package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// registerKeuanganRekonsiliasiRoutes — endpoint rekonsiliasi bank.
func (r *Router) registerKeuanganRekonsiliasiRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeTransaksiBasePath+"/rekonsiliasi", r.handleRekonsiliasiCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/rekonsiliasi/{id}", r.handleRekonsiliasiItem)
}

// =====================================================================
// Rekonsiliasi Collection — GET (list) | POST (jalankan baru)
// =====================================================================
func (r *Router) handleRekonsiliasiCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listRekonsiliasi(w, req, schoolID)
	case http.MethodPost:
		r.createRekonsiliasi(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listRekonsiliasi(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.RekonsiliasiFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		RekeningBankID:   req.URL.Query().Get("rekening_bank_id"),
		PeriodeTahun:     parseIntQuery(req, "periode_tahun", 0),
		Status:           req.URL.Query().Get("status"),
		Limit:            parseIntQuery(req, "limit", 20),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.rekonsiliasiService.ListRekonsiliasi(req.Context(), filter)
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

func (r *Router) createRekonsiliasi(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload models.CreateRekonsiliasiRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if payload.AcademicPeriodID == "" {
		payload.AcademicPeriodID = req.URL.Query().Get("academic_period_id")
	}
	if uid := req.Header.Get("X-User-ID"); uid != "" {
		payload.CreatedBy = uid
	}

	item, err := r.rekonsiliasiService.JalankanRekonsiliasi(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": item})
}

// =====================================================================
// Rekonsiliasi Item — GET (detail) | DELETE
// =====================================================================
func (r *Router) handleRekonsiliasiItem(w http.ResponseWriter, req *http.Request) {
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

	switch req.Method {
	case http.MethodGet:
		item, err := r.rekonsiliasiService.GetRekonsiliasi(req.Context(), schoolID, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})

	case http.MethodDelete:
		if err := r.rekonsiliasiService.DeleteRekonsiliasi(req.Context(), schoolID, id); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "rekonsiliasi dihapus"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}
