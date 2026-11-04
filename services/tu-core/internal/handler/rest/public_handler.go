package rest

import (
	"net/http"
	"sekolah-platform/services/tu-core/internal/service"
)

type PublicHandler struct {
	publicService *service.PublicService
}

func NewPublicHandler(publicService *service.PublicService) *PublicHandler {
	return &PublicHandler{publicService: publicService}
}

func (h *PublicHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	stats, err := h.publicService.GetDashboardStats(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data statistik"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total_siswa":  stats.TotalSiswa,
			"total_guru":   stats.TotalGuru,
			"total_rombel": stats.TotalRombel,
			"kompetensi_keahlian": 6,
		},
	})
}
