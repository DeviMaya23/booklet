package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewInternalAuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			got := c.Request().Header.Get("X-Booklet-Internal-Secret")
			if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			return next(c)
		}
	}
}
