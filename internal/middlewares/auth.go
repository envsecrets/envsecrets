package middlewares

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func WebhookHeader() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Hasura-Webhook-Secret",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("NHOST_WEBHOOK_SECRET"), nil
		},
	})
}
