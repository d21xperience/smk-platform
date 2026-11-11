package rest

import (
	"net/http"
	"sekolah-platform/services/tu-core/internal/service"
	"strconv"
)

type AgendaHandler struct {
	svc *service.AgendaService
}

func NewAgendaHandler(svc *service.AgendaService) *AgendaHandler {
	return &AgendaHandler{svc: svc}
}

func (h *AgendaHandler) List(w http.ResponseWriter, req *http.Request) {
	jenis := req.URL.Query().Get("jenis") // IN, OUT, BOTH
	bulan, _ := strconv.Atoi(req.URL.Query().Get("bulan"))
	tahun, _ := strconv.Atoi(req.URL.Query().Get("tahun"))

	if tahun == 0 {
		tahun = 2026 // Default
	}

	res, err := h.svc.ListAgenda(req.Context(), jenis, bulan, tahun)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *AgendaHandler) Export(w http.ResponseWriter, req *http.Request) {
	jenis := req.URL.Query().Get("jenis")
	bulan, _ := strconv.Atoi(req.URL.Query().Get("bulan"))
	tahun, _ := strconv.Atoi(req.URL.Query().Get("tahun"))

	fileURL, err := h.svc.GenerateExcel(req.Context(), jenis, bulan, tahun)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"file_url": fileURL})
}
