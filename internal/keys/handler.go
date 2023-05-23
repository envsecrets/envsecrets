package keys

import (
	"encoding/base64"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/keys/graphql"
	"github.com/labstack/echo/v4"
)

func GetPublicKey(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.GetPublicKeyOptions
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

	var key []byte
	if payload.Email != "" {
		result, err := graphql.GetPublicKeyByUserEmail(ctx, client, payload.Email)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Message: "Failed to get user's public key",
				Error:   err.Error(),
			})
		}
		key = result
	} else if payload.UserID != "" {
		result, err := graphql.GetPublicKeyByUserID(ctx, client, payload.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Message: "Failed to get user's public key",
				Error:   err.Error(),
			})
		}
		key = result
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully fetched user's public key",
		Data:    base64.StdEncoding.EncodeToString(key),
	})
}
