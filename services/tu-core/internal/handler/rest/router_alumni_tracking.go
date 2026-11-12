package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerAlumniTrackingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kesiswaan/tracer-study", r.wrapWithContext(r.handleAlumniTrackingList))
	mux.HandleFunc("/api/v1/kesiswaan/tracer-study/", r.wrapWithContext(r.handleAlumniTrackingUpdate))
}

func (r *Router) handleAlumniTrackingList(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	opCtx := middleware.GetOperationalContext(req.Context())
	list, err := r.alumniTrackingService.List(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

func (r *Router) handleAlumniTrackingUpdate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	opCtx := middleware.GetOperationalContext(req.Context())
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kesiswaan/tracer-study/")
	id := strings.TrimSpace(path)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	var reqBody service.UpdateAlumniTrackingRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	err := r.alumniTrackingService.Update(req.Context(), id, opCtx.SchoolID, opCtx.AcademicPeriodID, reqBody)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Data rekam jejak berhasil diperbarui"})
}
