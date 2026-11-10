package rest

import (
	"net/http"
)

func (r *Router) registerSarprasReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/sarpras/reports", r.wrapWithContext(r.handleSarprasReports))
}

func (r *Router) handleSarprasReports(w http.ResponseWriter, req *http.Request) {
	// Stub - Akan diimplementasikan di chat Sarpras
	writeJSON(w, http.StatusOK, map[string]string{"message": "Sarpras Reports endpoint"})
}
