package middleware

import (
	"net/http"
	"sekolah-platform/platform/context"

	"github.com/gin-gonic/gin"
)

func OperationalContextMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        schoolID := c.GetHeader("x-school-id")
        academicPeriodID := c.GetHeader("x-academic-period-id")
        userID := c.GetHeader("x-user-id")
        role := c.GetHeader("x-user-role")

        if schoolID == "" || academicPeriodID == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Missing operational context headers",
            })
            c.Abort()
            return
        }

        opCtx := &context.OperationalContext{
            SchoolID:         schoolID,
            AcademicPeriodID: academicPeriodID,
            UserID:           userID,
            Role:             role,
        }

        c.Set("operationalContext", opCtx)
        c.Next()
    }
}
