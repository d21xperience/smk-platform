package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// registerKeuanganRabRoutes — endpoint RAB & realisasi.
func (r *Router) registerKeuanganRabRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeTransaksiBasePath+"/rab", r.handleRabCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/rab/{id}", r.handleRabItem)
	mux.HandleFunc(financeTransaksiBasePath+"/rab/{id}/sahkan", r.handleRabSahkan)
	mux.HandleFunc(financeTransaksiBasePath+"/rab/tahun/{tahun}/realisasi", r.handleRabRealisasi)
}

// =====================================================================
// RAB Collection — GET (list) | POST (create)
// =====================================================================
func (r *Router) handleRabCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}

	switch req.Method {
	case http.MethodGet:
		r.listRab(w, req, schoolID)
	case http.MethodPost:
		r.createRab(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listRab(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.RabFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		TahunAnggaran:    parseIntQuery(req, "tahun_anggaran", 0),
		Status:           req.URL.Query().Get("status"),
		Limit:            parseIntQuery(req, "limit", 20),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.rabService.ListRab(req.Context(), filter)
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

func (r *Router) createRab(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload models.CreateRabRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID

	// Ambil user_id dari header (kalau ada), untuk disusun_oleh
	if uid := req.Header.Get("X-User-ID"); uid != "" {
		payload.DisusunOleh = uid
	}

	item, err := r.rabService.CreateRab(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": item})
}

// =====================================================================
// RAB Item — GET | PUT | DELETE
// =====================================================================
func (r *Router) handleRabItem(w http.ResponseWriter, req *http.Request) {
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
		item, err := r.rabService.GetRab(req.Context(), schoolID, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})

	case http.MethodPut:
		r.updateRab(w, req, schoolID, id)

	case http.MethodDelete:
		if err := r.rabService.DeleteRab(req.Context(), schoolID, id); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "RAB dihapus"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) updateRab(w http.ResponseWriter, req *http.Request, schoolID, id string) {
	var payload models.CreateRabRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	item, err := r.rabService.UpdateRab(req.Context(), id, payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

// =====================================================================
// RAB Sahkan — PUT (transisi DRAFT/DIAJUKAN → DISAHKAN)
// =====================================================================
func (r *Router) handleRabSahkan(w http.ResponseWriter, req *http.Request) {
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

	var payload models.SahkanRabRequest
	if req.ContentLength > 0 {
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
			return
		}
	}
	payload.SchoolID = schoolID
	payload.RabID = id
	if payload.DisahkanOleh == "" {
		payload.DisahkanOleh = req.Header.Get("X-User-ID")
	}

	item, err := r.rabService.SahkanRab(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}

// =====================================================================
// RAB Realisasi — GET (compute on-the-fly dari BKU)
// =====================================================================
func (r *Router) handleRabRealisasi(w http.ResponseWriter, req *http.Request) {
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
	tahun := parseIntQuery(req, "tahun", 0)
	if tahun == 0 {
		// Coba parse dari path value {tahun}
		tahun = parseIntQuery(req, "tahun_path", 0)
	}
	// Path value "tahun" tidak otomatis ke req — ambil dari URL
	// Go 1.22: PathValue("tahun")
	if tahun == 0 {
		if v := req.PathValue("tahun"); v != "" {
			tahun = parseIntQueryString(v)
		}
	}
	if tahun == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tahun wajib diisi"})
		return
	}

	sampaiBulan := parseIntQuery(req, "sampai_bulan", 12)
	report, err := r.rabService.GetRealisasi(req.Context(), schoolID, periodID, tahun, sampaiBulan)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": report})
}

// parseIntQueryString — parse string langsung ke int (0 kalau invalid)
func parseIntQueryString(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
