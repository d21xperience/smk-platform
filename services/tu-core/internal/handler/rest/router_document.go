package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerDocumentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/sarpras/documents", r.wrapWithContext(r.handleDocuments))
	mux.HandleFunc("/api/v1/sarpras/documents/", r.wrapWithContext(r.handleDocumentByID))
}

func (r *Router) handleDocuments(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		entityType := req.URL.Query().Get("entity_type")
		entityID := req.URL.Query().Get("entity_id")
		if entityType == "" || entityID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "entity_type dan entity_id wajib diisi"})
			return
		}
		data, err := r.documentService.GetDocuments(req.Context(), entityType, entityID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if data == nil { data = make([]*repository.InventoryDocument, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		if err := req.ParseMultipartForm(10 << 20); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal parse form"})
			return
		}
		file, handler, err := req.FormFile("file")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File wajib diunggah"})
			return
		}
		defer file.Close()

		docDate, _ := time.Parse("2006-01-02", req.FormValue("document_date"))
		doc := &repository.InventoryDocument{
			DocumentType:   req.FormValue("document_type"),
			Title:          req.FormValue("title"),
			Description:    req.FormValue("description"),
			DocumentNumber: req.FormValue("document_number"),
			DocumentDate:   &docDate,
			Visibility:     req.FormValue("visibility"),
		}

		opCtx := middleware.GetOperationalContext(req.Context())
		userID := "system"
		if opCtx != nil { userID = opCtx.UserID }

		entityType := req.FormValue("entity_type")
		entityID := req.FormValue("entity_id")

		if err := r.documentService.UploadAndLink(req.Context(), doc, entityType, entityID, file, handler, userID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Dokumen berhasil diunggah"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleDocumentByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/documents/")
	parts := splitPath(path)
	if len(parts) < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	id := parts[0]
	action := ""
	if len(parts) > 1 { action = parts[1] }

	opCtx := middleware.GetOperationalContext(req.Context())
	userID := "system"
	if opCtx != nil { userID = opCtx.UserID }

	switch action {
	case "verify":
		if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
		if err := r.documentService.VerifyDocument(req.Context(), id, userID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Dokumen diverifikasi"})

	case "reject":
		if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
		var body struct{ Reason string `json:"reason"` }
		json.NewDecoder(req.Body).Decode(&body)
		if err := r.documentService.RejectDocument(req.Context(), id, userID, body.Reason); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Dokumen ditolak"})

	case "download":
		if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
		url, err := r.documentService.GetDownloadURL(req.Context(), id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return
		}
		writeJSON(w, http.StatusOK, map[string]string{"url": url})

	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}
