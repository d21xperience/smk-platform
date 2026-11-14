package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// Register: Kode Rekening (COA) & Rekening Bank
// =====================================================================
func (r *Router) registerKeuanganMasterCOARoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeMasterBasePath+"/kode-rekening", r.handleKodeRekeningCollection)
	mux.HandleFunc(financeMasterBasePath+"/kode-rekening/{id}", r.handleKodeRekeningItem)
	mux.HandleFunc(financeMasterBasePath+"/rekening-bank", r.handleRekeningBankCollection)
	mux.HandleFunc(financeMasterBasePath+"/rekening-bank/{id}", r.handleRekeningBankItem)
}

// =====================================================================
// KODE REKENING
// =====================================================================
func (r *Router) handleKodeRekeningCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listKodeRekening(w, req, schoolID)
	case http.MethodPost:
		r.createKodeRekening(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listKodeRekening(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.KodeRekeningFilter{
		SchoolID: schoolID,
		Tipe:     req.URL.Query().Get("tipe"),
		Search:   req.URL.Query().Get("search"),
		IsAktif:  parseBoolPtrQuery(req, "is_aktif"),
		Limit:    parseIntQuery(req, "limit", 100),
		Offset:   parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.keuanganMasterService.ListKodeRekening(req.Context(), filter)
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

func (r *Router) createKodeRekening(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload models.KodeRekening
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if err := r.keuanganMasterService.CreateKodeRekening(req.Context(), &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": payload})
}

func (r *Router) handleKodeRekeningItem(w http.ResponseWriter, req *http.Request) {
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
		item, err := r.keuanganMasterService.GetKodeRekening(req.Context(), schoolID, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
	case http.MethodPut:
		var payload models.KodeRekening
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
			return
		}
		payload.ID = id
		payload.SchoolID = schoolID
		if err := r.keuanganMasterService.UpdateKodeRekening(req.Context(), &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": payload})
	case http.MethodDelete:
		if err := r.keuanganMasterService.DeleteKodeRekening(req.Context(), schoolID, id); err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "kode rekening dihapus"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// =====================================================================
// REKENING BANK
// =====================================================================
func (r *Router) handleRekeningBankCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listRekeningBank(w, req, schoolID)
	case http.MethodPost:
		r.createRekeningBank(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listRekeningBank(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.RekeningBankFilter{
		SchoolID: schoolID,
		Search:   req.URL.Query().Get("search"),
		IsAktif:  parseBoolPtrQuery(req, "is_aktif"),
		Limit:    parseIntQuery(req, "limit", 50),
		Offset:   parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.keuanganMasterService.ListRekeningBank(req.Context(), filter)
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

func (r *Router) createRekeningBank(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload models.RekeningBank
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if err := r.keuanganMasterService.CreateRekeningBank(req.Context(), &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": payload})
}

func (r *Router) handleRekeningBankItem(w http.ResponseWriter, req *http.Request) {
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
		item, err := r.keuanganMasterService.GetRekeningBank(req.Context(), schoolID, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
	case http.MethodPut:
		var payload models.RekeningBank
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
			return
		}
		payload.ID = id
		payload.SchoolID = schoolID
		if err := r.keuanganMasterService.UpdateRekeningBank(req.Context(), &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": payload})
	case http.MethodDelete:
		if err := r.keuanganMasterService.DeleteRekeningBank(req.Context(), schoolID, id); err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "rekening bank dihapus"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}
