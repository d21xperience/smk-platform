package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// registerKeuanganPenutupanRoutes — endpoint penutupan buku kas.
func (r *Router) registerKeuanganPenutupanRoutes(mux *http.ServeMux) {
	// Preview didaftarkan sebelum {id} agar ServeMux Go 1.22 memilih
	// pattern paling spesifik: /penutupan/preview > /penutupan/{id}
	mux.HandleFunc(financeTransaksiBasePath+"/penutupan/preview", r.handlePenutupanPreview)
	mux.HandleFunc(financeTransaksiBasePath+"/penutupan", r.handlePenutupanCollection)
	mux.HandleFunc(financeTransaksiBasePath+"/penutupan/{id}", r.handlePenutupanItem)
}

// =====================================================================
// Preview — GET snapshot Bku + Posisi sebelum user confirm tutup buku
//   query: ?academic_period_id=&periode_bulan=&periode_tahun=
// =====================================================================
func (r *Router) handlePenutupanPreview(w http.ResponseWriter, req *http.Request) {
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
	bulan := parseIntQuery(req, "periode_bulan", 0)
	tahun := parseIntQuery(req, "periode_tahun", 0)
	if bulan < 1 || bulan > 12 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "periode_bulan harus 1-12"})
		return
	}
	if tahun < 2000 || tahun > 2100 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "periode_tahun tidak valid"})
		return
	}

	snapshotBku, err := r.penutupanService.GetSnapshotBku(req.Context(), schoolID, periodID, bulan, tahun)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	snapshotPosisi, err := r.penutupanService.GetSnapshotPosisi(req.Context(), schoolID, periodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"periode_bulan":   bulan,
			"periode_tahun":   tahun,
			"snapshot_bku":    snapshotBku,
			"snapshot_posisi": snapshotPosisi,
		},
	})
}

// =====================================================================
// Collection — GET (list) | POST (tutup buku)
// =====================================================================
func (r *Router) handlePenutupanCollection(w http.ResponseWriter, req *http.Request) {
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		r.listPenutupan(w, req, schoolID)
	case http.MethodPost:
		r.createPenutupan(w, req, schoolID)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (r *Router) listPenutupan(w http.ResponseWriter, req *http.Request, schoolID string) {
	filter := models.PenutupanFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: req.URL.Query().Get("academic_period_id"),
		PeriodeTahun:     parseIntQuery(req, "periode_tahun", 0),
		Status:           req.URL.Query().Get("status"),
		Limit:            parseIntQuery(req, "limit", 20),
		Offset:           parseIntQuery(req, "offset", 0),
	}
	items, total, err := r.penutupanService.ListPenutupan(req.Context(), filter)
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

func (r *Router) createPenutupan(w http.ResponseWriter, req *http.Request, schoolID string) {
	var payload models.CreatePenutupanRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.SchoolID = schoolID
	if payload.AcademicPeriodID == "" {
		payload.AcademicPeriodID = req.URL.Query().Get("academic_period_id")
	}
	if uid := req.Header.Get("X-User-ID"); uid != "" {
		payload.ClosedBy = uid
	}

	item, err := r.penutupanService.TutupBuku(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": item})
}

// =====================================================================
// Item — GET (detail by ID)
// =====================================================================
func (r *Router) handlePenutupanItem(w http.ResponseWriter, req *http.Request) {
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
	item, err := r.penutupanService.GetPenutupan(req.Context(), schoolID, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": item})
}
