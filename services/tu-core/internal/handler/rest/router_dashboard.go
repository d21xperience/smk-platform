package rest

import (
	"net/http"

	"sekolah-platform/services/tu-core/internal/middleware"
)

func (r *Router) registerDashboardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/dashboard/summary", r.wrapWithContext(r.handleDashboardSummary))
}

func (r *Router) handleDashboardSummary(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil || opCtx.AcademicPeriodID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Header x-academic-period-id wajib diisi"})
		return
	}

	data, err := r.dashboardService.GetSummary(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}
