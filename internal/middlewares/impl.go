package middlewares

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/tokens"
	"github.com/envsecrets/envsecrets/utils"
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
	JWT_SIGNING_KEY, err := utils.GetJWTSecret()
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

			//	Decode the token
			payload, err := hex.DecodeString(key)
			if err != nil {
				return false, err
			}

			//	Initialize a new default context
			ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

			//	Initialize Hasura client with admin privileges
			client := clients.NewGQLClient(&clients.GQLConfig{
				Type: clients.HasuraClientType,
				Headers: []clients.Header{
					clients.XHasuraAdminSecretHeader,
				},
			})

			//	Hash the token to fetch it from database.
			hash := utils.SHA256Hash(payload)

			//	Decrypt the token.
			token, err := tokens.GetService().GetByHash(ctx, client, hash)
			if err != nil {
				return false, err
			}

			//	If the token is expired, return false.
			if token.IsExpired() {
				return false, echo.ErrUnauthorized
			}

			c.Set("token", token)

			return true, nil
		},
	})
}
