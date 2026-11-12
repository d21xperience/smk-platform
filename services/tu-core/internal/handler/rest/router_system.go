package rest

import (
	"net/http"
	"strings"
)

func (r *Router) registerSystemRoutes(mux *http.ServeMux) {
	// Auth & Public
	mux.HandleFunc("/api/v1/auth/login", r.authHandler.Login)
	mux.HandleFunc("/api/v1/public/stats", r.publicHandler.GetStats)

	// Context & Preferences
	mux.HandleFunc("/api/v1/context/periods", r.contextHandler.GetPeriods)
	mux.HandleFunc("/api/v1/context/toggle-psb", r.contextHandler.TogglePSB)

	mux.HandleFunc("/api/v1/users/me/preferences", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.contextHandler.GetPreference(w, req)
		case http.MethodPut:
			r.contextHandler.UpdatePreference(w, req)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})

	// Academic Periods
	mux.HandleFunc("/api/v1/academic/periods", r.handleCreateAcademicPeriod)
	mux.HandleFunc("/api/v1/academic/periods/", r.handleAcademicPeriodsByID)

	// User Management
	mux.HandleFunc("/api/v1/users", r.handleUsers)
	mux.HandleFunc("/api/v1/users/", r.handleUsersByID)

	// School Profile
	mux.HandleFunc("/api/v1/school", r.handleSchool)

	// Legacy Stubs (Tetap dipertahankan agar rute lama tidak 404)
	mux.HandleFunc("/api/v1/mutations", r.stubHandler)
	mux.HandleFunc("/api/v1/assessments", r.stubHandler)
	mux.HandleFunc("/api/v1/finance", r.stubHandler)

}

func (r *Router) handleSchool(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodGet:
		r.schoolHandler.GetSchool(w, req)
	case http.MethodPut:
		r.schoolHandler.UpdateSchool(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleUsers(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodGet:
		r.userHandler.ListUsers(w, req)
	case http.MethodPost:
		r.userHandler.CreateUser(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleUsersByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodPut:
		r.userHandler.UpdateUser(w, req)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleCreateAcademicPeriod(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	r.contextHandler.CreatePeriod(w, req.WithContext(ctx))
}

func (r *Router) handleAcademicPeriodsByID(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context headers"})
		return
	}
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodPut:
		if strings.HasSuffix(req.URL.Path, "/activate") {
			r.contextHandler.ActivatePeriod(w, req)
			return
		}
	case http.MethodDelete:
		r.contextHandler.DeletePeriod(w, req)
		return
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
}

func (r *Router) stubHandler(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "Not implemented"})
}
