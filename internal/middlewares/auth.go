package middlewares

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func WebhookHeader() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-webhook-secret",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("NHOST_WEBHOOK_SECRET"), nil
		},
	})
}

func HasWebhookHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Response().Header().Get("x-webhook-secret") != os.Getenv("NHOST_WEBHOOK_SECRET") {
			return c.String(http.StatusUnauthorized, "invalid webhook auth")
		}
		return nil
	}
}
