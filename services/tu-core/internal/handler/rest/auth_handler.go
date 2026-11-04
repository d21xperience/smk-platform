package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format request tidak valid"})
		return
	}

	result, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"token": result.Token,
			"user": map[string]interface{}{
				"id":                 result.User.UserID,
				"name":               result.User.Name,
				"username":           result.User.Name,
				"role":               result.User.Role,
				"school_id":          result.User.SchoolID,
				"school_name":        result.User.SchoolName, // <-- Tambahan
				"academic_period_id": result.AcademicPeriodID,
				"academic_year":      result.AcademicYear,
				"semester":           result.Semester,
				"permissions":        []string{"*"},
			},
		},
	})
}
