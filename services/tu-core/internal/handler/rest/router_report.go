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
	semesterID := "20261"
	if opCtx != nil && opCtx.AcademicPeriodID != "" {
		semesterID = opCtx.AcademicPeriodID
	}

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
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format start_date tidak valid (YYYY-MM-DD)"})
		return
	}
	
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format end_date tidak valid (YYYY-MM-DD)"})
		return
	}

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
