package auth

import (
	"encoding/base64"
	"net/http"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/auth/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func SignupHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SignupOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
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

	//	Call the service handler.
	if err := Signup(ctx, client, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to register the user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "Verification email sent to your inbox!",
		Data:    payload.Email,
	})
}

func UpdatePasswordHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.UpdatePasswordOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Extract the user's email from JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.Claims)

	//	Check whether user has keys.
	_, err := keys.GetByUserID(ctx, client, claims.Hasura.UserID)
	if err != nil {
		apiError := clients.ParseExternal(err)

		//	If key pair does not exist for this user,
		//	issue them a new key pair.
		if apiError.IsType(clients.ErrorTypeDoesNotExist) {

			//	Generate Key pair
			pair, err := keys.GenerateKeyPair(payload.NewPassword)
			if err != nil {
				return c.JSON(http.StatusBadRequest, &clients.APIResponse{
					Message: "Failed to issue a fresh key pair",
					Error:   err.Error(),
				})
			}

			//	Upload the keys to their cloud account.
			if err := keys.Create(ctx, client, &keyCommons.CreateOptions{
				PublicKey:    base64.StdEncoding.EncodeToString(pair.PublicKey),
				PrivateKey:   base64.StdEncoding.EncodeToString(pair.PrivateKey),
				ProtectedKey: base64.StdEncoding.EncodeToString(pair.ProtectedKey),
				Salt:         base64.StdEncoding.EncodeToString(pair.Salt),
			}); err != nil {
				return c.JSON(http.StatusBadRequest, &clients.APIResponse{
					Message: "Failed to issue a fresh key pair",
					Error:   err.Error(),
				})
			}
		} else {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Error: err.Error(),
			})
		}
	}

	//	TODO: If a key-pair already exists,
	//	update the user's protection key.

	//	Initialize HTTP client
	httpClient := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service handler.
	if err := UpdatePassword(ctx, httpClient, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to update the password",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "password successfuly updated",
	})
}
