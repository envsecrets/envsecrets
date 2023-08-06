package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/internal/actions"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/invites"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/middlewares"
	"github.com/envsecrets/envsecrets/internal/payments"
	"github.com/envsecrets/envsecrets/internal/secrets"
	"github.com/envsecrets/envsecrets/internal/tokens"
	"github.com/envsecrets/envsecrets/internal/triggers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := godotenv.Load(".env.development")
	if err != nil {
		log.Println("Error loading .env.development file")
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	//	Add healthcheck endpoint
	e.GET("/healthz", healthz)

	//	Initialize the routes to skip JWT auth
	skipRoutes := []string{
		"/triggers",
		"integrations",
		"healthz",
		"invites",
		"/secrets",
		"/payments/server/webhook",
		"/auth/signup",
		"/auth/validate-password",
	}

	skipper := func(c echo.Context) bool {
		for _, item := range skipRoutes {
			if strings.Contains(c.Request().URL.Path, item) {
				return true
			}
		}
		return false
	}
	e.Use(middlewares.JWTAuth(skipper))

	//	API	Version 1 Group
	v1Group := e.Group("/v1")

	//	Hasura triggers group
	triggers.AddRoutes(v1Group)

	//	Hasura actions group
	actions.AddRoutes(v1Group)

	//	Authentication group
	auth.AddRoutes(v1Group)

	//	Keys group
	keys.AddRoutes(v1Group)

	//	Secrets group
	secrets.AddRoutes(v1Group)

	//	Integrations group
	integrations.AddRoutes(v1Group)

	//	Environments group
	environments.AddRoutes(v1Group)

	//	Invites group
	invites.AddRoutes(v1Group)

	//	Payments group
	payments.AddRoutes(v1Group)

	//	Tokens group
	tokens.AddRoutes(v1Group)

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

// Healthcheck endpoint
func healthz(c echo.Context) error {
	return c.String(http.StatusOK, "API is healthy")
}
