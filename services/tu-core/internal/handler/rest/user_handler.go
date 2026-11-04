package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func isAdmin(req *http.Request) bool {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		return false
	}
	role := strings.ToLower(strings.TrimSpace(opCtx.Role))
	return role == "admin"
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(req) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak."})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	users, err := h.userSvc.ListUsers(req.Context(), opCtx.SchoolID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(req) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak."})
		return
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	if err := h.userSvc.CreateUser(req.Context(), opCtx.SchoolID, body.Username, body.Password, body.Role); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Pengguna berhasil ditambahkan"})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	if !isAdmin(req) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Akses ditolak."})
		return
	}

	// Path: /api/v1/users/{id}/toggle atau /api/v1/users/{id}/reset-password
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/users/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
		return
	}
	userID := parts[0]

	// Aksi 1: Toggle Status
	if len(parts) > 1 && parts[1] == "toggle" {
		var body struct {
			IsActive bool `json:"is_active"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
			return
		}
		if err := h.userSvc.ToggleStatus(req.Context(), userID, body.IsActive); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Status berhasil diubah"})
		return
	}

	// Aksi 2: Reset Password (BARU)
	if len(parts) > 1 && parts[1] == "reset-password" {
		var body struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
			return
		}
		if err := h.userSvc.ResetPassword(req.Context(), userID, body.Password); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Password berhasil direset"})
		return
	}

	// Aksi 3: Update Role
	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}
	if err := h.userSvc.UpdateRole(req.Context(), userID, body.Role); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil diubah"})
}
