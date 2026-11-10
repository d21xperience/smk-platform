package rest

import (
	"net/http"

	"sekolah-platform/services/tu-core/internal/middleware"
)

func (r *Router) registerStatisticRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kepegawaian/statistics/dashboard", r.wrapWithContext(r.handleStatisticDashboard))
}

func (r *Router) handleStatisticDashboard(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	startDate := req.URL.Query().Get("start_date")
	endDate := req.URL.Query().Get("end_date")

	data, err := r.ptkStatisticService.GetDashboardData(req.Context(), opCtx.SchoolID, startDate, endDate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, data)
}
