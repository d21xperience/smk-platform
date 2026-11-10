package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerMaintenanceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/sarpras/damage-reports", r.wrapWithContext(r.handleDamageReports))
	mux.HandleFunc("/api/v1/sarpras/damage-reports/", r.wrapWithContext(r.handleDamageReportByID))
	mux.HandleFunc("/api/v1/sarpras/maintenance-schedules", r.wrapWithContext(r.handleMaintenanceSchedules))
	mux.HandleFunc("/api/v1/sarpras/maintenance-schedules/", r.wrapWithContext(r.handleMaintenanceScheduleByID))
	mux.HandleFunc("/api/v1/sarpras/maintenance-logs", r.wrapWithContext(r.handleMaintenanceLogs))
}

// ... (handleDamageReports, handleDamageReportByID, handleMaintenanceSchedules, handleMaintenanceScheduleByID tetap sama)

func (r *Router) handleDamageReports(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	semesterID := ""
	userID := "system"
	if opCtx != nil {
		semesterID = opCtx.AcademicPeriodID
		userID = opCtx.UserID
	}

	switch req.Method {
	case http.MethodGet:
		data, err := r.maintenanceService.GetAllDamageReports(req.Context(), semesterID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if data == nil {
			data = make([]*repository.DamageReport, 0)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		var body repository.DamageReport
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		body.SemesterID = semesterID
		body.ReporterID = userID
		res, err := r.maintenanceService.CreateDamageReport(req.Context(), &body)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "Tiket berhasil dibuat", "data": res})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleDamageReportByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/damage-reports/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	if req.Method == http.MethodPut {
		var body struct {
			Status string  `json:"status"`
			Vendor string  `json:"vendor"`
			Cost   float64 `json:"cost"`
			Notes  string  `json:"notes"`
		}
		json.NewDecoder(req.Body).Decode(&body)
		if err := r.maintenanceService.UpdateDamageReport(req.Context(), id, body.Status, body.Vendor, body.Cost, body.Notes); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Tiket berhasil diupdate"})
		return
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
}

func (r *Router) handleMaintenanceSchedules(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	semesterID := ""
	if opCtx != nil {
		semesterID = opCtx.AcademicPeriodID
	}

	switch req.Method {
	case http.MethodGet:
		data, err := r.maintenanceService.GetAllSchedules(req.Context(), semesterID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if data == nil {
			data = make([]*repository.MaintenanceSchedule, 0)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		var body struct {
			AssetName       string `json:"asset_name"`
			InventoryItemID string `json:"inventory_item_id"`
			ScopeLocation   string `json:"scope_location"`
			Frequency       string `json:"frequency"`
			NextServiceDate string `json:"next_service_date"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}

		nextDate, _ := time.Parse("2006-01-02", body.NextServiceDate)
		sch := &repository.MaintenanceSchedule{
			AssetName: body.AssetName, InventoryItemID: body.InventoryItemID, ScopeLocation: body.ScopeLocation,
			Frequency: body.Frequency, NextServiceDate: nextDate, SemesterID: semesterID,
		}
		res, err := r.maintenanceService.CreateSchedule(req.Context(), sch)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "Jadwal berhasil dibuat", "data": res})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleMaintenanceScheduleByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/maintenance-schedules/")
	id := strings.TrimSuffix(path, "/complete")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	if strings.HasSuffix(path, "/complete") && req.Method == http.MethodPost {
		opCtx := middleware.GetOperationalContext(req.Context())
		userID := "system"
		if opCtx != nil {
			userID = opCtx.UserID
		}

		var body struct {
			ServiceDate string  `json:"service_date"`
			Cost        float64 `json:"cost"`
			Notes       string  `json:"notes"`
		}
		json.NewDecoder(req.Body).Decode(&body)
		sDate, _ := time.Parse("2006-01-02", body.ServiceDate)

		if err := r.maintenanceService.CompleteSchedule(req.Context(), id, sDate, body.Cost, body.Notes, userID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Servis berhasil diselesaikan"})
		return
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
}

func (r *Router) handleMaintenanceLogs(w http.ResponseWriter, req *http.Request) {
	startDate := req.URL.Query().Get("start_date")
	endDate := req.URL.Query().Get("end_date")
	isPDF := req.URL.Query().Get("format") == "pdf"

	if isPDF {
		pdfBuf, err := r.maintenanceService.GenerateMaintenanceBookPDF(req.Context(), startDate, endDate)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"Buku_Perawatan.pdf\"")
		w.WriteHeader(http.StatusOK)
		w.Write(pdfBuf.Bytes())
		return
	}

	data, err := r.maintenanceService.GetMaintenanceLogs(req.Context(), startDate, endDate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if data == nil {
		data = make([]*repository.MaintenanceLog, 0)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}
