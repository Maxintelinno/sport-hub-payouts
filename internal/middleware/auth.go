package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware is a placeholder for JWT authentication
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing Authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid Authorization format"})
		}

		// In a real implementation, you would verify the JWT here
		token := parts[1]
		
		// Placeholder: Extract owner_id from token
		// For now, if token is "mock-owner-123", set owner_id to "123"
		// In production, this should be the sub claim from the verified JWT
		ownerID := token 
		if ownerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		c.Set("owner_id", ownerID)
		return next(c)
	}
}
