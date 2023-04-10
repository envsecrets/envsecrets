package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/internal/actions"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/invites"
	"github.com/envsecrets/envsecrets/internal/payments"
	"github.com/envsecrets/envsecrets/internal/secrets"
	"github.com/envsecrets/envsecrets/internal/triggers"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt"
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

	//	Load the JWT signing from env vars
	JWT_SIGNING_KEY := getJWTSecret("NHOST_JWT_SECRET")

	//	Initialize the routes to skip JWT auth
	skipRoutes := []string{
		"/triggers",
		"integrations",
		"healthz",
		"invites",
		"/payments/server/webhook",
	}

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(JWT_SIGNING_KEY.Key),
		SigningMethod: JWT_SIGNING_KEY.Type,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.Claims)
		},

		//	Define function to skip certain routes from JWT auth
		Skipper: func(c echo.Context) bool {
			for _, item := range skipRoutes {
				if strings.Contains(c.Request().URL.Path, item) {
					return true
				}
			}
			return false
		},
	}))

	//	API	Version 1 Group
	v1Group := e.Group("/v1")

	//	Hasura triggers group
	triggers.AddRoutes(v1Group)

	//	Hasura actions group
	actions.AddRoutes(v1Group)

	//	Secrets group
	secrets.AddRoutes(v1Group)

	//	Integrations group
	integrations.AddRoutes(v1Group)

	//	Invites group
	invites.AddRoutes(v1Group)

	//	Payments group
	payments.AddRoutes(v1Group)

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

type JWTSecret struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

//	Get the JWT key from NHOST variables
func getJWTSecret(variable string) *JWTSecret {
	var response JWTSecret
	payload := os.Getenv(variable)
	json.Unmarshal([]byte(payload), &response)
	return &response
}

//	Healthcheck endpoint
func healthz(c echo.Context) error {
	return c.String(http.StatusOK, "API is healthy")
}
