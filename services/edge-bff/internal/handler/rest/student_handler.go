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
	studentID := c.Param("id")

	student, err := h.repo.FindStudentByID(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(500, errorResponse("INTERNAL_ERROR", err.Error()))
		return
	}
	if student == nil {
		c.JSON(404, errorResponse("NOT_FOUND", "Student tidak ditemukan"))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    student,
		"error":   nil,
	})
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
