package middlewares

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/envsecrets/envsecrets/cli/auth"
	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/tokens"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func WebhookHeader() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:" + string(clients.HasuraWebhookSecret),
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("NHOST_WEBHOOK_SECRET"), nil
		},
	})
}

func JWTAuth(skipper middleware.Skipper) echo.MiddlewareFunc {

	//	Load the JWT signing from env vars
	JWT_SIGNING_KEY, err := globalCommons.GetJWTSecret()
	if err != nil {
		log.Fatal("unable to get the JWT secret: ", err)
	}

	config := echojwt.Config{
		SigningKey:    []byte(JWT_SIGNING_KEY.Key),
		SigningMethod: JWT_SIGNING_KEY.Type,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.Claims)
		},
	}

	if skipper != nil {
		config.Skipper = skipper
	}

	return echojwt.WithConfig(config)

}

func TokenHeader() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().Header.Get(string(clients.AuthorizationHeader)) != ""
		},
		KeyLookup: "header:" + string(clients.TokenHeader),
		Validator: func(key string, c echo.Context) (bool, error) {

			//	Initialize a new default context
			ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

			//	Initialize Hasura client with admin privileges
			client := clients.NewGQLClient(&clients.GQLConfig{
				Type: clients.HasuraClientType,
				Headers: []clients.Header{
					clients.XHasuraAdminSecretHeader,
				},
			})

			//	Decode the token
			payload, err := hex.DecodeString(key)
			if err != nil {
				return false, err
			}

			//	Decrypt the token.
			token, err := tokens.GetService().Decrypt(ctx, client, payload)
			if err != nil {
				return false, err
			}

			//	Set the environment ID in echo's context.
			c.Set("env_id", token.EnvID)
			c.Set("org_key", token.OrgKey)

			return true, nil
		},
	})
}
