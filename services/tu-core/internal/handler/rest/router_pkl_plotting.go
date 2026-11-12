package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerPklPlottingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kesiswaan/pkl-plotting", r.wrapWithContext(r.handlePklPlottingListCreate))
	mux.HandleFunc("/api/v1/kesiswaan/pkl-plotting/", r.wrapWithContext(r.handlePklPlottingByID))
}

func (r *Router) handlePklPlottingListCreate(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())

	switch req.Method {
	case http.MethodGet:
		list, err := r.pklPlottingService.List(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})

	case http.MethodPost:
		var reqBody service.PklPlottingRequest
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body: " + err.Error()})
			return
		}
		if reqBody.Nama == "" || reqBody.Perusahaan == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama siswa dan Perusahaan wajib diisi"})
			return
		}
		err := r.pklPlottingService.Create(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID, reqBody)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Plotting berhasil dibuat"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handlePklPlottingByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kesiswaan/pkl-plotting/")
	id := strings.TrimSpace(path)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	switch req.Method {
	case http.MethodPut:
		var reqBody service.PklPlottingRequest
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		err := r.pklPlottingService.Update(req.Context(), id, opCtx.SchoolID, opCtx.AcademicPeriodID, reqBody)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Plotting berhasil diperbarui"})

	case http.MethodDelete:
		err := r.pklPlottingService.Delete(req.Context(), id, opCtx.SchoolID, opCtx.AcademicPeriodID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Plotting berhasil dihapus"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
