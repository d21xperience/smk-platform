package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerGoodsReceiptRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/sarpras/receipts", r.wrapWithContext(r.handleGoodsReceipts))
	mux.HandleFunc("/api/v1/sarpras/receipts/", r.wrapWithContext(r.handleGoodsReceiptByID))
	mux.HandleFunc("/api/v1/sarpras/receipts/upload", r.wrapWithContext(r.handleReceiptUpload))
}

func (r *Router) handleGoodsReceipts(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	semesterID := ""
	if opCtx != nil { semesterID = opCtx.AcademicPeriodID }

	switch req.Method {
	case http.MethodGet:
		data, err := r.goodsReceiptService.GetAllReceipts(req.Context(), semesterID)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		if data == nil { data = make([]*repository.GoodsReceipt, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		var body repository.GoodsReceipt
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return
		}
		if opCtx != nil {
			body.ReceivedBy = opCtx.UserID
			body.SemesterID = opCtx.AcademicPeriodID
		}
		res, err := r.goodsReceiptService.CreateReceipt(req.Context(), &body)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "Penerimaan berhasil dibuat: " + res.BastNumber, "data": res})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleGoodsReceiptByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/receipts/")
	parts := splitPath(path)
	if len(parts) < 1 { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }
	id := parts[0]

	opCtx := middleware.GetOperationalContext(req.Context())
	userID := "system"
	if opCtx != nil { userID = opCtx.UserID }

	if len(parts) == 1 && req.Method == http.MethodGet {
		data, err := r.goodsReceiptService.GetAllReceipts(req.Context(), "") // Simplified for dummy
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, data)
		return
	}

	if len(parts) == 2 {
		action := parts[1]
		if action == "verify" && req.Method == http.MethodPost {
			if err := r.goodsReceiptService.VerifyReceipt(req.Context(), id, userID); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Penerimaan berhasil diverifikasi"})
			return
		}
		if action == "reject" && req.Method == http.MethodPost {
			var body struct{ Reason string `json:"reason"` }
			json.NewDecoder(req.Body).Decode(&body)
			if err := r.goodsReceiptService.RejectReceipt(req.Context(), id, userID, body.Reason); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Penerimaan ditolak"})
			return
		}
		if action == "pdf" && req.Method == http.MethodGet {
			pdfBuf, err := r.goodsReceiptService.GeneratePDF(req.Context(), id)
			if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"BAST_%s.pdf\"", id))
			w.WriteHeader(http.StatusOK)
			w.Write(pdfBuf.Bytes())
			return
		}
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method or path not allowed"})
}

func (r *Router) handleReceiptUpload(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	err := req.ParseMultipartForm(10 << 20)
	if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal parse form"}); return }
	file, handler, err := req.FormFile("document")
	if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File tidak ditemukan"}); return }
	defer file.Close()

	objectName := fmt.Sprintf("bast_doc_%d_%s", time.Now().Unix(), handler.Filename)
	url, err := r.minioClient.UploadFile(req.Context(), "documents", objectName, file, handler)
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}
