package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig adalah konfigurasi JWT middleware.
type JWTConfig struct {
	Secret string
	Issuer string
}

// JWTMiddleware memvalidasi JWT token.
func JWTMiddleware(config JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "UNAUTHORIZED",
					"message": "Missing Authorization header",
				},
			})
			c.Abort()
			return
		}

		// Extract token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(401, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INVALID_TOKEN",
					"message": "Invalid Authorization header format",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["sub"])
			c.Set("user_role", claims["role"])
			c.Set("school_id", claims["school_id"])
		}

		c.Next()
	}
}
