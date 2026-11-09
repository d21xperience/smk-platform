package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerProcurementRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/procurement/requests", r.wrapWithContext(r.handleProcurementRequests))
	mux.HandleFunc("/api/v1/procurement/requests/", r.wrapWithContext(r.handleProcurementRequestByID))
}

func (r *Router) handleProcurementRequests(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	semesterID := ""
	if opCtx != nil {
		semesterID = opCtx.AcademicPeriodID
	}

	switch req.Method {
	case http.MethodGet:
		data, err := r.procurementService.GetAllRequests(req.Context(), semesterID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if data == nil {
			data = make([]*repository.ProcurementRequest, 0)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		var body repository.ProcurementRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
			return
		}

		if opCtx != nil {
			body.RequesterID = opCtx.UserID
			body.SemesterID = opCtx.AcademicPeriodID
		}

		res, err := r.procurementService.CreateRequest(req.Context(), &body)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Pengajuan berhasil dibuat dengan nomor: " + res.RequestNumber,
			"data":    res,
		})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleProcurementRequestByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/procurement/requests/")
	if path == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	parts := splitPath(path)
	id := parts[0]

	// Jika path hanya ID (untuk PUT edit / DELETE)
	if len(parts) == 1 {
		switch req.Method {
		case http.MethodPut:
			var body repository.ProcurementRequest
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
				return
			}
			body.ID = id
			if err := r.procurementService.UpdateRequest(req.Context(), &body); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Pengajuan berhasil diupdate"})

		case http.MethodDelete:
			if err := r.procurementService.DeleteRequest(req.Context(), id); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Pengajuan berhasil dihapus"})

		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
		return
	}

	// Jika path ID/status (untuk update status)
	if len(parts) == 2 && parts[1] == "status" {
		if req.Method != http.MethodPut {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}

		var body struct {
			Status          string `json:"status"`
			RejectionReason string `json:"rejection_reason"`
			Notes           string `json:"notes"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
			return
		}

		opCtx := middleware.GetOperationalContext(req.Context())
		approverID := "system"
		if opCtx != nil {
			approverID = opCtx.UserID
		}

		var err error
		switch body.Status {
		case "APPROVED":
			err = r.procurementService.ApproveRequest(req.Context(), id, approverID, body.Notes)
		case "REJECTED":
			err = r.procurementService.RejectRequest(req.Context(), id, approverID, body.RejectionReason, body.Notes)
		case "PROCURED":
			err = r.procurementService.MarkAsProcured(req.Context(), id)
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Status tidak valid"})
			return
		}

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "Status berhasil diupdate"})
		return
	}

	writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
}
