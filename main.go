package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/envsecrets/envsecrets/internal/events"
	"github.com/joho/godotenv"
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

	/* 	JWT_SIGNING_KEY := getJWTSecret("NHOST_JWT_SECRET")
	   	e.Use(echojwt.WithConfig(echojwt.Config{
	   		SigningKey:    JWT_SIGNING_KEY.Key,
	   		SigningMethod: JWT_SIGNING_KEY.Type,
	   	}))
	*/

	//	TODO: Add webhook validation middleware

	//	API	Version 1 Group
	v1Group := e.Group("/api/v1")

	//	Hasura events group
	events.AddRoutes(v1Group)

	//	Secrets group
	//secrets := v1Group.Group("/secrets")

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
