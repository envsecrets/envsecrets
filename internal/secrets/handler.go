package secrets

import (
	"encoding/base64"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/labstack/echo/v4"
)

func SetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SetRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   echo.HeaderAuthorization,
				Value: c.Request().Header.Get(echo.HeaderAuthorization),
			},
		},
	})

	//	Base64 encode the secret value
	base64Value := base64.StdEncoding.EncodeToString([]byte(payload.Secret.Value.(string)))
	payload.Secret.Value = base64Value

	//	Call the service function.
	if err := Set(ctx, client, &commons.SetOptions{
		Path:   payload.Path,
		Secret: payload.Secret,
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
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   echo.HeaderAuthorization,
				Value: c.Request().Header.Get(echo.HeaderAuthorization),
			},
		},
	})

	//	Call the service function.
	secret, err := Get(ctx, client, &payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to set the secret",
			Error:   err.Message,
		})
	}

	//	Base64 decode the secret value
	data, er := base64.StdEncoding.DecodeString(secret.Value.(string))
	if er != nil {
		return c.JSON(http.StatusBadRequest, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to set the secret",
			Error:   err.Message,
		})
	}

	secret.Value = string(data)

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully got the secret",
		Data:    secret,
	})
}
