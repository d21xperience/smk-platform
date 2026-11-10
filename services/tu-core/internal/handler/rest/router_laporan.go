package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerLaporanRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kepegawaian/laporan", r.wrapWithContext(r.handleLaporanList))
	mux.HandleFunc("/api/v1/kepegawaian/laporan/generate", r.wrapWithContext(r.handleGenerateLaporan))
	mux.HandleFunc("/api/v1/kepegawaian/laporan/", r.wrapWithContext(r.handleLaporanByID))
}

func (r *Router) handleLaporanList(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	list, err := r.laporanBerkalaRepo.List(req.Context(), opCtx.SchoolID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handleGenerateLaporan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	var body struct {
		Tahun    int `json:"tahun"`
		Triwulan int `json:"triwulan"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	var start, end time.Time
	year := body.Tahun
	switch body.Triwulan {
	case 1: start = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC); end = time.Date(year, 3, 31, 23, 59, 59, 0, time.UTC)
	case 2: start = time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC); end = time.Date(year, 6, 30, 23, 59, 59, 0, time.UTC)
	case 3: start = time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC); end = time.Date(year, 9, 30, 23, 59, 59, 0, time.UTC)
	case 4: start = time.Date(year, 10, 1, 0, 0, 0, 0, time.UTC); end = time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Triwulan tidak valid"})
		return
	}

	jsonData, err := r.laporanAggregatorService.AggregateData(req.Context(), opCtx.SchoolID, start, end)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	laporan := &repository.LaporanBerkala{
		SchoolID: opCtx.SchoolID, Tahun: year, Triwulan: body.Triwulan,
		TanggalMulai: start, TanggalSelesai: end, DataJSON: jsonData,
	}

	if err := r.laporanBerkalaRepo.CreateDraft(req.Context(), laporan); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_ = r.laporanBerkalaRepo.LogTimeline(req.Context(), laporan.ID, "", "DRAFT", opCtx.UserID, "Draft digenerate manual")

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Draft laporan berhasil dibuat",
		"data":    laporan,
	})
}

func (r *Router) handleLaporanByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/laporan/")
	parts := splitPath(path)
	if len(parts) < 2 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL"})
		return
	}
	id := parts[0]
	action := parts[1]

	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var body struct {
		Notes string `json:"notes"`
		Email string `json:"email"`
	}
	if req.Body != http.NoBody {
		_ = json.NewDecoder(req.Body).Decode(&body)
	}

	var err error
	switch action {
	case "generate-documents":
		err = r.laporanBerkalaService.GenerateDocuments(req.Context(), id, opCtx.SchoolID)
	case "review":
		err = r.laporanBerkalaService.ReviewLaporan(req.Context(), id, opCtx.UserID, body.Notes)
	case "approve":
		err = r.laporanBerkalaService.ApproveLaporan(req.Context(), id, opCtx.UserID, body.Notes)
	case "send":
		err = r.laporanBerkalaService.SendLaporan(req.Context(), id, opCtx.UserID, r.laporanEmailService, body.Email)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unknown action"})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Aksi " + action + " berhasil"})
}
