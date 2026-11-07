package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

type CorrespondenceHandler struct {
	svc *service.CorrespondenceService
}

func NewCorrespondenceHandler(svc *service.CorrespondenceService) *CorrespondenceHandler {
	return &CorrespondenceHandler{svc: svc}
}

// Create
func (h *CorrespondenceHandler) Create(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	var body struct {
		Type            string                 `json:"type"`
		JenisSurat      string                 `json:"jenis_surat"`
		Kategori        string                 `json:"kategori"`
		KodeKlasifikasi string                 `json:"kode_klasifikasi"`
		Number          string                 `json:"number"`
		Date            string                 `json:"date"`
		Subject         string                 `json:"subject"`
		From            string                 `json:"from"`
		To              string                 `json:"to"`
		Description     string                 `json:"description"`
		AttachmentURL   string                 `json:"attachment_url"`
		FileURL         string                 `json:"file_url"`
		TTDURL          string                 `json:"ttd_url"`
		Konten          map[string]interface{} `json:"konten"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format tanggal harus YYYY-MM-DD"})
		return
	}

	svcReq := service.CreateCorrespondenceRequest{
		Type:            body.Type,
		JenisSurat:      body.JenisSurat,
		Kategori:        body.Kategori,
		KodeKlasifikasi: body.KodeKlasifikasi,
		Number:          body.Number,
		Date:            date,
		Subject:         body.Subject,
		From:            body.From,
		To:              body.To,
		Description:     body.Description,
		AttachmentURL:   body.AttachmentURL,
		FileURL:         body.FileURL,
		TTDURL:          body.TTDURL,
		Konten:          body.Konten,
	}

	res, err := h.svc.CreateCorrespondence(req.Context(), svcReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, res)
}

// List
func (h *CorrespondenceHandler) List(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	q := req.URL.Query()
	filter := repository.CorrespondenceFilter{
		Type:            q.Get("type"),
		Status:          q.Get("status"),
		JenisSurat:      q.Get("jenis_surat"),
		Kategori:        q.Get("kategori"),
		KodeKlasifikasi: q.Get("kode_klasifikasi"),
		DateFrom:        q.Get("date_from"),
		DateTo:          q.Get("date_to"),
		Search:          q.Get("search"),
	}

	res, err := h.svc.ListCorrespondences(req.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// GetByID
func (h *CorrespondenceHandler) GetByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	id := parts[0]

	res, err := h.svc.GetCorrespondence(req.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// Update
func (h *CorrespondenceHandler) Update(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	id := parts[0]

	var body struct {
		Kategori        string                 `json:"kategori"`
		KodeKlasifikasi string                 `json:"kode_klasifikasi"`
		Subject         string                 `json:"subject"`
		Description     string                 `json:"description"`
		AttachmentURL   string                 `json:"attachment_url"`
		FileURL         string                 `json:"file_url"`
		TTDURL          string                 `json:"ttd_url"`
		Konten          map[string]interface{} `json:"konten"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	svcReq := service.UpdateCorrespondenceRequest{
		Kategori:        body.Kategori,
		KodeKlasifikasi: body.KodeKlasifikasi,
		Subject:         body.Subject,
		Description:     body.Description,
		AttachmentURL:   body.AttachmentURL,
		FileURL:         body.FileURL,
		TTDURL:          body.TTDURL,
		Konten:          body.Konten,
	}

	res, err := h.svc.UpdateCorrespondence(req.Context(), id, svcReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// Delete
func (h *CorrespondenceHandler) Delete(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	id := parts[0]

	if err := h.svc.DeleteCorrespondence(req.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Surat berhasil dihapus"})
}

// Archive
func (h *CorrespondenceHandler) Archive(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "archive" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	id := parts[0]

	res, err := h.svc.ArchiveCorrespondence(req.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// DownloadPDF memicu generate PDF (jika belum ada) dan menyajikan file-nya
func (h *CorrespondenceHandler) DownloadPDF(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "pdf" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	id := parts[0]

	fileURL, err := h.svc.GeneratePDF(req.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Ambil path file dari URL (asumsi disimpan di ./storage/pdfs)
	filePath := fmt.Sprintf("./storage/pdfs/%s.pdf", id)

	// Set header untuk download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"Surat_%s.pdf\"", id))

	// Sajikan file
	http.ServeFile(w, req, filePath)

	_ = fileURL // Biarkan ini di sini untuk menandakan bahwa URL telah di-update di DB
}
