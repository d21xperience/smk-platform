package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerBorrowingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/sarpras/borrowings", r.wrapWithContext(r.handleBorrowings))
	mux.HandleFunc("/api/v1/sarpras/borrowings/", r.wrapWithContext(r.handleBorrowingByID))
}

func (r *Router) handleBorrowings(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	semesterID := ""
	userID := "system"
	userName := "System"
	if opCtx != nil {
		semesterID = opCtx.AcademicPeriodID
		userID = opCtx.UserID
	}
	
	switch req.Method {
	case http.MethodGet:
		data, err := r.borrowingService.GetAll(req.Context(), semesterID)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		if data == nil { data = make([]*repository.Borrowing, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		var body struct {
			InventoryItemID    string `json:"inventory_item_id"`
			BorrowerName       string `json:"borrower_name"`
			PurposeType        string `json:"purpose_type"`
			PurposeDescription string `json:"purpose_description"`
			BorrowDate         string `json:"borrow_date"`
			ExpectedReturnDate string `json:"expected_return_date"`
			UsageLocation      string `json:"usage_location"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		
		borrowDate, _ := time.Parse("2006-01-02", body.BorrowDate)
		returnDate, _ := time.Parse("2006-01-02", body.ExpectedReturnDate)
		
		// Fallback name if from body
		if body.BorrowerName != "" { userName = body.BorrowerName }
		
		b := &repository.Borrowing{
			InventoryItemID:    body.InventoryItemID,
			PurposeType:        body.PurposeType,
			PurposeDescription: body.PurposeDescription,
			BorrowDate:         borrowDate,
			ExpectedReturnDate: returnDate,
			UsageLocation:      body.UsageLocation,
		}
		
		res, err := r.borrowingService.CreateBorrowing(req.Context(), b, userID, userName)
		if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "Peminjaman berhasil diajukan: " + res.BorrowingNumber, "data": res})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleBorrowingByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/borrowings/")
	parts := splitPath(path)
	if len(parts) < 1 { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }
	id := parts[0]

	opCtx := middleware.GetOperationalContext(req.Context())
	userID := "system"
	if opCtx != nil { userID = opCtx.UserID }

	if len(parts) == 1 && req.Method == http.MethodGet {
		data, err := r.borrowingService.GetByID(req.Context(), id)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, data)
		return
	}

	if len(parts) == 2 {
		action := parts[1]
		
		// Fetch latest data
		b, err := r.borrowingService.GetByID(req.Context(), id)
		if err != nil { writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()}); return }
		
		switch action {
		case "approve":
			if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
			var body struct {
				ApprovalStep int    `json:"approval_step"`
				Decision     string `json:"decision"`
				Notes        string `json:"notes"`
			}
			json.NewDecoder(req.Body).Decode(&body)
			approval := &repository.BorrowingApproval{
				BorrowingID:  id,
				ApprovalStep: body.ApprovalStep,
				ApproverID:   userID,
				Decision:     body.Decision,
				Notes:        body.Notes,
			}
			err = r.borrowingService.Approve(req.Context(), b, approval)
			
		case "pickup":
			if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
			var body struct {
				Condition string `json:"condition"`
				Photos    string `json:"photos"`
				Notes     string `json:"notes"`
			}
			json.NewDecoder(req.Body).Decode(&body)
			err = r.borrowingService.Pickup(req.Context(), b, userID, body.Condition, body.Photos, body.Notes)
			
		case "return":
			if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
			var body struct {
				Condition string `json:"condition"`
				Photos    string `json:"photos"`
				Notes     string `json:"notes"`
			}
			json.NewDecoder(req.Body).Decode(&body)
			err = r.borrowingService.Return(req.Context(), b, userID, body.Condition, body.Photos, body.Notes)
		}

		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Aksi berhasil diproses"})
		return
	}

	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method or path not allowed"})
}
