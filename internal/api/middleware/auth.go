package middleware

import (
	"net/http"
	"strings"

	"github.com/jexlor/votingapp/internal/auth"
	"github.com/labstack/echo/v4"
)

func JWTAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			// Try Authorization header first
			authHeader := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// If no header, try cookie
			if tokenString == "" {
				cookie, err := c.Cookie("auth_token")
				if err == nil {
					tokenString = cookie.Value
				}
			}

			// Still empty? Unauthorized
			if tokenString == "" {
				c.Redirect(http.StatusSeeOther, "/v1/login")
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			// Parse and validate
			claims, err := auth.ParseJWT(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			// Attach user ID to context
			c.Set("userID", claims.UserID)
			return next(c)
		}
	}
}
