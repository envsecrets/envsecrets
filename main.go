package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/secrets"
	"github.com/envsecrets/envsecrets/internal/triggers"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//	Load the JWT signing from env vars
	JWT_SIGNING_KEY := getJWTSecret("NHOST_JWT_SECRET")

	//	Initialize the routes to skip JWT auth
	skipRoutes := []string{
		"/triggers",
		"integrations",
	}

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(JWT_SIGNING_KEY.Key),
		SigningMethod: JWT_SIGNING_KEY.Type,

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

	//	Secrets group
	secrets.AddRoutes(v1Group)

	//	Integrations group
	integrations.AddRoutes(v1Group)

	/* 	// Routes
	   	secrets.POST("/set", api.SetSecret)
	   	secrets.GET("/get", api.GetSecret)
	   	secrets.GET("/get/versions", api.GetSecretVersions)
	   	secrets.GET("/list", api.ListSecrets)

	   	//	Invites group
	   	invites := v1Group.Group("/invites")
	   	invites.POST("/:id/accept", api.AcceptInvite)
	*/
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
