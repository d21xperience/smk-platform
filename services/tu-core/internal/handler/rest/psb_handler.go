package rest

import (
	"net/http"

	"sekolah-platform/services/tu-core/internal/middleware"
)

func (r *Router) handlePSBCalonSiswa(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	opCtx := middleware.GetOperationalContext(ctx)
	list, err := r.psbService.GetCalonSiswa(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

func (r *Router) handlePSBGenerateKelas(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	opCtx := middleware.GetOperationalContext(ctx)
	if err := r.psbService.GenerateKelas(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Kelas berhasil di-generate"})
}
