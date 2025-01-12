package middleware

import (
	"net/http"
	"strings"

	"gc-buku/utils"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "missing authorization header",
			})
		}

		token := strings.Replace(auth, "Bearer ", "", 1)
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
		}

		c.Set("user_id", claims.UserID)
		return next(c)
	}
}
