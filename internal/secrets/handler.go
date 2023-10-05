package secrets

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	tokenCommons "github.com/envsecrets/envsecrets/internal/tokens/commons"
	"github.com/labstack/echo/v4"
)

func SetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	secret, err := Set(ctx, client, &commons.SetOptions{
		EnvID: payload.EnvID,
		Data:  payload.Data,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to set the secrets",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully set the secret",
		Data:    secret,
	})
}

func DeleteHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.DeleteRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	if _, err := Delete(ctx, client, &commons.DeleteSecretOptions{
		EnvID:   payload.EnvID,
		Key:     payload.Key,
		Version: payload.Version,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to delete the secret",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully delete the secret",
	})
}

func GetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.GetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
	})

	//	If the user has passed an authorization header,
	//	use that in GraphQL client.
	//	Else if they are authenticating using a token,
	//	it is safe to use the admin token.
	if c.Request().Header.Get(echo.HeaderAuthorization) != "" {
		client.Authorization = c.Request().Header.Get(echo.HeaderAuthorization)
	} else if c.Request().Header.Get(string(clients.TokenHeader)) != "" {
		client.Headers = append(client.Headers, clients.XHasuraAdminSecretHeader)
	} else {
		return echo.ErrUnauthorized
	}

	var token *tokenCommons.Token
	if c.Get("token") != nil {
		token = c.Get("token").(*tokenCommons.Token)
	}

	//	Override the env_id set by token middleware.
	if token != nil {
		payload.EnvID = token.EnvID
	}

	//	Call the service function.
	secret, err := Get(ctx, client, &commons.GetOptions{
		Key:     payload.Key,
		EnvID:   payload.EnvID,
		Version: payload.Version,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get the secret",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully got the secret",
		Data: commons.GetResponse{
			Secret: secret,
			Token:  token,
		},
	})
}
