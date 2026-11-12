package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerAlumniLegalisirRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kesiswaan/legalisir", r.wrapWithContext(r.handleAlumniLegalisirListCreate))
	mux.HandleFunc("/api/v1/kesiswaan/legalisir/", r.wrapWithContext(r.handleAlumniLegalisirByID))
}

func (r *Router) handleAlumniLegalisirListCreate(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())

	switch req.Method {
	case http.MethodGet:
		list, err := r.alumniLegalisirService.List(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})

	case http.MethodPost:
		var reqBody service.CreateAlumniLegalisirRequest
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body: " + err.Error()})
			return
		}
		if reqBody.Nama == "" || reqBody.NoIjazah == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama dan No Ijazah wajib diisi"})
			return
		}
		err := r.alumniLegalisirService.Create(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID, reqBody)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Berhasil dibuat"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleAlumniLegalisirByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kesiswaan/legalisir/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	id := parts[0]

	if len(parts) == 2 && parts[1] == "status" && req.Method == http.MethodPut {
		var reqBody struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		err := r.alumniLegalisirService.UpdateStatus(req.Context(), id, opCtx.SchoolID, opCtx.AcademicPeriodID, reqBody.Status)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Status diperbarui"})
		return
	}

	if len(parts) == 1 && req.Method == http.MethodDelete {
		err := r.alumniLegalisirService.Delete(req.Context(), id, opCtx.SchoolID, opCtx.AcademicPeriodID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Data dihapus"})
		return
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Endpoint not found"})
}
