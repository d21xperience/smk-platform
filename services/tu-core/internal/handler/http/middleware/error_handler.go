package middleware

import (
	"errors"
	"net/http"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch {
			case errors.Is(err, repository.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{
					"code":    "NOT_FOUND",
					"message": err.Error(),
				})
			case errors.Is(err, repository.ErrDuplicateEntry):
				c.JSON(http.StatusConflict, gin.H{
					"code":    "DUPLICATE_ENTRY",
					"message": err.Error(),
				})
			case errors.Is(err, repository.ErrInvalidInput):
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "INVALID_INPUT",
					"message": err.Error(),
				})
			case errors.Is(err, service.ErrMissingOperationalContext):
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "MISSING_CONTEXT",
					"message": err.Error(),
				})
			case errors.Is(err, service.ErrUnauthorized):
				c.JSON(http.StatusForbidden, gin.H{
					"code":    "UNAUTHORIZED",
					"message": err.Error(),
				})
			case errors.Is(err, service.ErrBusinessRuleViolation):
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"code":    "BUSINESS_RULE_VIOLATION",
					"message": err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Terjadi kesalahan internal",
				})
			}
		}
	}
}
