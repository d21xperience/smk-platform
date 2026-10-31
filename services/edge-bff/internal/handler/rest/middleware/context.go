// Package middleware menyediakan HTTP middleware untuk edge-bff.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	platformctx "sekolah-platform/platform/context"
)

const (
	HeaderSchoolID       = "x-school-id"
	HeaderSchoolName     = "x-school-name"
	HeaderAcademicYear   = "x-academic-year"
	HeaderSemester       = "x-semester"
	HeaderAcademicPeriod = "x-academic-period-id"
	HeaderUserID         = "x-user-id"
	HeaderUserRole       = "x-user-role"
)

type contextKey string

const operationalContextKey contextKey = "operational_context"

// OperationalContextMiddleware mengekstrak OperationalContext dari HTTP headers.
func OperationalContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		schoolID := c.GetHeader(HeaderSchoolID)
		schoolName := c.GetHeader(HeaderSchoolName)
		academicYear := c.GetHeader(HeaderAcademicYear)
		periodID := c.GetHeader(HeaderAcademicPeriod)
		if periodID != "" {
			// Validasi format: 4 digit tahun + 1 digit semester
			if len(periodID) == 5 {
				// Simpan ke dalam Context Gin
				c.Set("academic_period_id", periodID)
				c.Set("academic_year", periodID[:4])
				c.Set("semester", periodID[4:])
			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Format x-academic-period-id tidak valid"})
				return
			}
		}

		userID := c.GetHeader(HeaderUserID)
		role := c.GetHeader(HeaderUserRole)

		// Parse semester dari periodID jika academicYear tidak ada
		semester := 0
		if periodID != "" && academicYear == "" {
			var s int
			academicYear, s = platformctx.ParsePeriodID(periodID)
			semester = s
		}

		opCtx := &platformctx.OperationalContext{
			SchoolID:         schoolID,
			SchoolName:       schoolName,
			AcademicYear:     academicYear,
			Semester:         semester,
			AcademicPeriodID: periodID,
			UserID:           userID,
			Role:             role,
		}

		// Derive periodID jika belum ada
		if opCtx.AcademicPeriodID == "" && opCtx.AcademicYear != "" && opCtx.Semester > 0 {
			opCtx.AcademicPeriodID = opCtx.DerivePeriodID()
		}

		// Validate
		if err := opCtx.Validate(); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"data":    nil,
				"error": map[string]interface{}{
					"code":    "INVALID_CONTEXT",
					"message": err.Error(),
				},
			})
			c.Abort()
			return
		}

		c.Set(string(operationalContextKey), opCtx)
		c.Next()
	}
}

// GetOperationalContext mengambil OperationalContext dari gin.Context.
func GetOperationalContext(c *gin.Context) *platformctx.OperationalContext {
	if v, exists := c.Get(string(operationalContextKey)); exists {
		if opCtx, ok := v.(*platformctx.OperationalContext); ok {
			return opCtx
		}
	}
	return nil
}
