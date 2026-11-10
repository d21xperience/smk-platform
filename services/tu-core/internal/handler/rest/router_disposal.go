package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerDisposalRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/sarpras/disposal", r.wrapWithContext(r.handleDisposals))
	mux.HandleFunc("/api/v1/sarpras/disposal/", r.wrapWithContext(r.handleDisposalByID))
}

func (r *Router) handleDisposals(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		opCtx := middleware.GetOperationalContext(req.Context())
		semesterID := "20261"
		if opCtx != nil { semesterID = opCtx.AcademicPeriodID }
		
		data, err := r.disposalService.GetAll(req.Context(), semesterID)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		if data == nil { data = make([]*repository.DisposalRequest, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		var body struct {
			ReasonCategory string                       `json:"reason_category"`
			ProposedMethod string                       `json:"proposed_method"`
			Items          []repository.DisposalItem    `json:"items"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }

		opCtx := middleware.GetOperationalContext(req.Context())
		userID := "system"
		semesterID := "20261"
		if opCtx != nil { userID = opCtx.UserID; semesterID = opCtx.AcademicPeriodID }

		disposalReq := &repository.DisposalRequest{
			ReasonCategory: body.ReasonCategory,
			ProposedMethod: body.ProposedMethod,
			ProposalDate:   time.Now(),
			Status:         "DRAFT",
		}

		res, err := r.disposalService.CreateRequest(req.Context(), disposalReq, body.Items, userID, semesterID)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "Usulan berhasil dibuat", "data": res})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleDisposalByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/disposal/")
	parts := splitPath(path)
	if len(parts) < 1 { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }
	id := parts[0]

	if req.Method == http.MethodGet && len(parts) == 1 {
		data, err := r.disposalService.GetByID(req.Context(), id)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, data)
		return
	}

	if len(parts) == 2 {
		action := parts[1]
		var err error

		switch action {
		case "verify":
			if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
			reqData, _ := r.disposalService.GetByID(req.Context(), id)
			var body repository.DisposalVerification
			json.NewDecoder(req.Body).Decode(&body)
			err = r.disposalService.Verify(req.Context(), reqData, &body)

		case "approve":
			if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
			reqData, _ := r.disposalService.GetByID(req.Context(), id)
			var body repository.DisposalApproval
			json.NewDecoder(req.Body).Decode(&body)
			err = r.disposalService.Approve(req.Context(), reqData, &body)

		case "execute":
			if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
			reqData, _ := r.disposalService.GetByID(req.Context(), id)
			var body struct {
				ExecutionDate     string `json:"execution_date"`
				ExecutionMethod   string `json:"execution_method"`
				ExecutionLocation string `json:"execution_location"`
			}
			json.NewDecoder(req.Body).Decode(&body)
			reqData.ExecutionMethod = body.ExecutionMethod
			reqData.ExecutionLocation = body.ExecutionLocation
			err = r.disposalService.Execute(req.Context(), reqData)
		}

		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Aksi berhasil diproses"})
		return
	}

	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method or path not allowed"})
}
