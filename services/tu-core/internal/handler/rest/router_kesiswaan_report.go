package rest

import (
	"fmt"
	"net/http"

	"sekolah-platform/services/tu-core/internal/middleware"
)

// registerKesiswaanReportRoutes mendaftarkan route untuk Pelaporan Kesiswaan.
func (r *Router) registerKesiswaanReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kesiswaan/pelaporan/statistik", r.wrapWithContext(r.handleKesiswaanReportStatistik))
	mux.HandleFunc("/api/v1/kesiswaan/pelaporan/export/excel", r.wrapWithContext(r.handleKesiswaanReportExcel))
	mux.HandleFunc("/api/v1/kesiswaan/pelaporan/export/pdf", r.wrapWithContext(r.handleKesiswaanReportPDF))
}

func (r *Router) handleKesiswaanReportStatistik(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	opCtx := middleware.GetOperationalContext(req.Context())
	data, err := r.kesiswaanReportService.GetStatistik(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (r *Router) handleKesiswaanReportExcel(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	opCtx := middleware.GetOperationalContext(req.Context())
	buf, filename, err := r.kesiswaanReportService.GenerateExcel(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (r *Router) handleKesiswaanReportPDF(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	opCtx := middleware.GetOperationalContext(req.Context())
	buf, filename, err := r.kesiswaanReportService.GeneratePDF(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
