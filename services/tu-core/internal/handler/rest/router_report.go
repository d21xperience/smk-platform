package rest

import (
	"fmt"
	"net/http"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
)

func (r *Router) registerReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/sarpras/dashboard", r.wrapWithContext(r.handleDashboardStats))
	mux.HandleFunc("/api/v1/sarpras/reports/comprehensive", r.wrapWithContext(r.handleComprehensiveReport))
}

func (r *Router) handleDashboardStats(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	
	opCtx := middleware.GetOperationalContext(req.Context())
	semesterID := ""
	if opCtx != nil { semesterID = opCtx.AcademicPeriodID }

	data, err := r.sarprasReportService.GetDashboardStats(req.Context(), semesterID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

func (r *Router) handleComprehensiveReport(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	
	startDateStr := req.URL.Query().Get("start_date")
	endDateStr := req.URL.Query().Get("end_date")
	
	if startDateStr == "" || endDateStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "start_date dan end_date wajib diisi"})
		return
	}
	
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	pdfBuf, err := r.sarprasReportService.GenerateComprehensivePDF(req.Context(), startDate, endDate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"Laporan_Komprehensif_%s.pdf\"", startDateStr))
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBuf.Bytes())
}
