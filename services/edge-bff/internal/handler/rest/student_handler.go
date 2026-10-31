// Package rest menyediakan HTTP handlers untuk edge-bff.
package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sekolah-platform/services/edge-bff/internal/handler/rest/middleware"
	"sekolah-platform/services/edge-bff/internal/repository"
)

// StudentHandler adalah HTTP handler untuk student endpoints.
type StudentHandler struct {
	repo *repository.SQLiteRepository
}

// NewStudentHandler membuat StudentHandler baru.
func NewStudentHandler(repo *repository.SQLiteRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

// RegisterRoutes mendaftarkan routes untuk student.
func (h *StudentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	students := rg.Group("/students")
	{
		students.GET("", h.ListStudents)
		students.GET("/:id", h.GetStudent)
	}
}

// ListStudents mengembalikan daftar siswa.
func (h *StudentHandler) ListStudents(c *gin.Context) {
	opCtx := middleware.GetOperationalContext(c)
	if opCtx == nil {
		c.JSON(400, errorResponse("INVALID_CONTEXT", "Operational context missing"))
		return
	}

	// Query parameters
	status := c.Query("status")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	students, total, err := h.repo.ListStudents(c.Request.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID, status, search, page, limit)
	if err != nil {
		c.JSON(500, errorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"items":      students,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": totalPages,
		},
		"error": nil,
	})
}

// GetStudent mengembalikan detail siswa.
func (h *StudentHandler) GetStudent(c *gin.Context) {
	// 1. Ambil SchoolID dari Middleware
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Operational context missing: school_id"})
		return
	}

	studentID := c.Param("id")

	// 2. Gunakan method scoped
	student, err := h.repo.FindStudentByIDAndSchoolID(c.Request.Context(), studentID, schoolID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// Lakukan hal yang sama untuk fungsi UpdateStudentStatus di handler
func (h *StudentHandler) UpdateStudentStatus(c *gin.Context) {
	schoolID := c.GetString("school_id")
	if schoolID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Operational context missing"})
		return
	}

	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.UpdateStudentStatus(c.Request.Context(), id, schoolID, req.Status)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to update: cross-tenant or not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// errorResponse membuat response error yang konsisten.
func errorResponse(code, message string) gin.H {
	return gin.H{
		"success": false,
		"data":    nil,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

// successResponse membuat response sukses.
func successResponse(data interface{}) gin.H {
	return gin.H{
		"success": true,
		"data":    data,
		"error":   nil,
	}
}

// Placeholder untuk http.StatusOK (digunakan di beberapa tempat)
var _ = http.StatusOK
