package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

type AttendanceHandler struct {
	attSvc *service.AttendanceService
}

func NewAttendanceHandler(attSvc *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attSvc: attSvc}
}

// CreateSession (POST /api/v1/attendance/sessions)
func (h *AttendanceHandler) CreateSession(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	var body struct {
		ClassID     string `json:"class_id"`
		SubjectID   string `json:"subject_id"`
		TeacherID   string `json:"teacher_id"`
		SessionDate string `json:"session_date"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
		Notes       string `json:"notes"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	// Parsing waktu (asumsi format: YYYY-MM-DD dan HH:MM)
	date, err := time.Parse("2006-01-02", body.SessionDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format tanggal harus YYYY-MM-DD"})
		return
	}
	start, _ := time.Parse("15:04", body.StartTime)
	end, _ := time.Parse("15:04", body.EndTime)

	reqData := service.CreateSessionRequest{
		ClassID:     body.ClassID,
		SubjectID:   body.SubjectID,
		TeacherID:   body.TeacherID,
		SessionDate: date,
		StartTime:   start,
		EndTime:     end,
		Notes:       body.Notes,
	}

	session, err := h.attSvc.CreateSession(req.Context(), reqData)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

// SubmitAttendance (POST /api/v1/attendance/sessions/{id}/submit)
func (h *AttendanceHandler) SubmitAttendance(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	// Ambil ID dari path /api/v1/attendance/sessions/{id}/submit
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/attendance/sessions/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "submit" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Path tidak valid"})
		return
	}
	sessionID := parts[0]

	var body struct {
		Records []struct {
			StudentID string `json:"student_id"`
			Status    string `json:"status"`
			Note      string `json:"note"`
		} `json:"records"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	var inputs []service.AttendanceRecordInput
	for _, r := range body.Records {
		inputs = append(inputs, service.AttendanceRecordInput{
			StudentID: r.StudentID,
			Status:    r.Status,
			Note:      r.Note,
		})
	}

	reqData := service.SubmitSessionAttendanceRequest{
		SessionID: sessionID,
		Records:   inputs,
	}

	session, err := h.attSvc.SubmitAttendance(req.Context(), reqData)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, session)
}

// ListSessions (GET /api/v1/attendance/sessions)
func (h *AttendanceHandler) ListSessions(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	classID := req.URL.Query().Get("class_id")
	startDateStr := req.URL.Query().Get("start_date")
	endDateStr := req.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" { startDate, _ = time.Parse("2006-01-02", startDateStr) }
	if endDateStr != "" { endDate, _ = time.Parse("2006-01-02", endDateStr) }

	sessions, err := h.attSvc.ListSessions(req.Context(), classID, "", startDate, endDate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, sessions)
}

// GetRecap (GET /api/v1/attendance/recap)
// Menghitung rekap mingguan/semester dari data sesi
func (h *AttendanceHandler) GetRecap(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	classID := req.URL.Query().Get("class_id")
	monthStr := req.URL.Query().Get("month")
	yearStr := req.URL.Query().Get("year")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	if month == 0 { month = int(time.Now().Month()) }
	if year == 0 { year = time.Now().Year() }

	// Memanggil fungsi GetAttendanceSummary yang sudah ada di Service
	summaries, err := h.attSvc.GetAttendanceSummary(req.Context(), classID, month, year)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, summaries)
}
