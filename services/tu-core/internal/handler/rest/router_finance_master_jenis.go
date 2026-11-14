package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

const financeMasterBasePath = "/api/v1/keuangan/master"

// =====================================================================
// Register: Jenis Iuran & Tarif Iuran
// =====================================================================
func (r *Router) registerKeuanganMasterJenisRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeMasterBasePath+"/jenis-iuran", r.handleJenisIuranCollection)
	mux.HandleFunc(financeMasterBasePath+"/jenis-iuran/{id}", r.handleJenisIuranItem)
	mux.HandleFunc(financeMasterBasePath+"/tarif-iuran", r.handleTarifIuranCollection)
	mux.HandleFunc(financeMasterBasePath+"/tarif-iuran/{id}", r.handleTarifIuranItem)
}

// =====================================================================
// JENIS IURAN
// =====================================================================
func (r *Router) handleJenisIuranCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi (query ?school_id= atau header X-School-ID)"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listJenisIuran(w, req, schoolID)
	case http.MethodPost:
		r.createJenisIuran(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listJenisIuran(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.JenisIuranFilter{
		SchoolID: schoolID,
		Kategori: req.URL.Query().Get("kategori"),
		Search:   req.URL.Query().Get("search"),
		IsAktif:  parseBoolPtrQuery(req, "is_aktif"),
		Limit:    parseIntQuery(req, "limit", 50),
		Offset:   parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.keuanganMasterService.ListJenisIuran(req.Context(), filter)
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

func (r *Router) createJenisIuran(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload models.JenisIuran
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if err := r.keuanganMasterService.CreateJenisIuran(req.Context(), &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": payload})
}

func (r *Router) handleJenisIuranItem(w http.ResponseWriter, req *http.Request) {
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
		item, err := r.keuanganMasterService.GetJenisIuran(req.Context(), schoolID, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
	case http.MethodPut:
		var payload models.JenisIuran
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
			return
		}
		payload.ID = id
		payload.SchoolID = schoolID
		if err := r.keuanganMasterService.UpdateJenisIuran(req.Context(), &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": payload})
	case http.MethodDelete:
		if err := r.keuanganMasterService.DeleteJenisIuran(req.Context(), schoolID, id); err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "jenis iuran dihapus"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// =====================================================================
// TARIF IURAN
// =====================================================================
func (r *Router) handleTarifIuranCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listTarifIuran(w, req, schoolID)
	case http.MethodPost:
		r.createTarifIuran(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listTarifIuran(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.TarifIuranFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		JenisIuranID:     req.URL.Query().Get("jenis_iuran_id"),
		Kelas:            req.URL.Query().Get("kelas"),
		Limit:            parseIntQuery(req, "limit", 50),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.keuanganMasterService.ListTarifIuran(req.Context(), filter)
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

func (r *Router) createTarifIuran(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload models.TarifIuran
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if err := r.keuanganMasterService.CreateTarifIuran(req.Context(), &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": payload})
}

func (r *Router) handleTarifIuranItem(w http.ResponseWriter, req *http.Request) {
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
		item, err := r.keuanganMasterService.GetTarifIuran(req.Context(), schoolID, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
	case http.MethodPut:
		var payload models.TarifIuran
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
			return
		}
		payload.ID = id
		payload.SchoolID = schoolID
		if err := r.keuanganMasterService.UpdateTarifIuran(req.Context(), &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": payload})
	case http.MethodDelete:
		if err := r.keuanganMasterService.DeleteTarifIuran(req.Context(), schoolID, id); err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "tarif iuran dihapus"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}
