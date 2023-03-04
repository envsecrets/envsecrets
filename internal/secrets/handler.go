package secrets

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/labstack/echo/v4"
)

func SetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	if err := Set(ctx, client, &commons.SetSecretOptions{
		KeyPath:    payload.OrgID,
		EnvID:      payload.EnvID,
		Data:       payload.Data,
		KeyVersion: payload.KeyVersion,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to set the secret",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully set the secret",
	})
}

func GetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.GetRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	secret, err := Get(ctx, client, &commons.GetSecretOptions{
		Key:     payload.Key,
		KeyPath: payload.OrgID,
		EnvID:   payload.EnvID,
		Version: payload.Version,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to get the secret",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully got the secret",
		Data:    secret.Data[payload.Key],
	})
}
