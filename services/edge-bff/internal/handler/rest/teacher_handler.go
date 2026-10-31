package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sekolah-platform/services/edge-bff/internal/handler/rest/middleware"
	"sekolah-platform/services/edge-bff/internal/repository"
)

// TeacherHandler adalah HTTP handler untuk teacher endpoints.
type TeacherHandler struct {
	repo *repository.SQLiteRepository
}

// NewTeacherHandler membuat TeacherHandler baru.
func NewTeacherHandler(repo *repository.SQLiteRepository) *TeacherHandler {
	return &TeacherHandler{repo: repo}
}

// RegisterRoutes mendaftarkan routes untuk teacher.
func (h *TeacherHandler) RegisterRoutes(rg *gin.RouterGroup) {
	teachers := rg.Group("/teachers")
	{
		teachers.GET("", h.ListTeachers)
		teachers.GET("/:id", h.GetTeacher)
	}
}

// ListTeachers mengembalikan daftar guru.
func (h *TeacherHandler) ListTeachers(c *gin.Context) {
	opCtx := middleware.GetOperationalContext(c)
	if opCtx == nil {
		c.JSON(400, errorResponse("INVALID_CONTEXT", "Operational context missing"))
		return
	}

	subject := c.Query("subject")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	teachers, total, err := h.repo.ListTeachers(c.Request.Context(), opCtx.SchoolID, subject, status, page, limit)
	if err != nil {
		c.JSON(500, errorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"items":      teachers,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": totalPages,
		},
		"error": nil,
	})
}

// GetTeacher mengembalikan detail guru.
func (h *TeacherHandler) GetTeacher(c *gin.Context) {
	// 1. Ambil SchoolID dari Middleware
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Operational context missing: school_id"})
		return
	}

	teacherID := c.Param("id")

	// 2. Gunakan method scoped
	teacher, err := h.repo.FindTeacherByIDAndSchoolID(c.Request.Context(), teacherID, schoolID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

// Placeholder
var _ = http.StatusOK
