package middlewares

import (
	"net/http"
	"strings"

	"gin-gonic/helpers"
	"gin-gonic/utils"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware middleware untuk validasi JWT token
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required", nil)
			c.Abort()
			return
		}

		// Check if token starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format. Use 'Bearer <token>'", nil)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate JWT token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Set user information in context for use in handlers
		c.Set("user_id", claims["user_id"])
		c.Set("user_email", claims["email"])
		c.Set("user_name", claims["name"])
		c.Set("user_role", claims["role"])

		c.Next()
	}
}

//middleware admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exits := c.Get("user_role")
		if !exits || role != "admin" {
			helpers.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Hanya admin dibenarkan mengakses ini")
			c.Abort()
			return
		}
	}
}

// OptionalJWTMiddleware middleware yang opsional (tidak wajib ada token)
func OptionalJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without user context
			c.Next()
			return
		}

		// Check if token starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Invalid format, continue without user context
			c.Next()
			return
		}

		tokenString := tokenParts[1]

		// Validate JWT token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			// Invalid token, continue without user context
			c.Next()
			return
		}

		// Set user information in context for use in handlers
		c.Set("user_id", claims["user_id"])
		c.Set("user_email", claims["email"])
		c.Set("user_name", claims["name"])

		c.Next()
	}
}
