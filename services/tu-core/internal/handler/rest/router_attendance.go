package rest

import (
	"net/http"
)

func (r *Router) registerAttendanceRoutes(mux *http.ServeMux) {
	// Rute untuk membuat sesi absensi (jika diperlukan di masa depan)
	mux.HandleFunc("/api/v1/attendance/sessions", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.attendanceHandler.CreateSession(w, req)
		case http.MethodGet:
			r.attendanceHandler.ListSessions(w, req)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	})

	// Rute untuk submit absensi (jika diperlukan)
	mux.HandleFunc("/api/v1/attendance/sessions/{id}/submit", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			r.attendanceHandler.SubmitAttendance(w, req)
		} else {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	})

	// Rute untuk rekap kehadiran (yang sedang kita butuhkan)
	mux.HandleFunc("/api/v1/attendance/recap", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			r.attendanceHandler.GetRecap(w, req)
		} else {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	})
}
