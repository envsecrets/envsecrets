package tokens

import (
	"net/http"
	"time"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/tokens/commons"
	"github.com/labstack/echo/v4"
)

func CreateHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.CreateRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with user's token
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Create the token
	expiry, er := time.ParseDuration(payload.Expiry)
	if er != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to parse expiry duration",
			Error:   er.Error(),
		})
	}

	token, err := Create(ctx, client, &commons.CreateServiceOptions{
		EnvID:  payload.EnvID,
		Expiry: expiry,
		Name:   payload.Name,
	})
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to create the token"),
			Error:   err.Error.Error(),
		})
	}

	//	Re-Initialize a new Hasura client with admin privileges
	client = clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Get the token hash using the new admin client
	token, err = Get(ctx, client, token.ID)
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to create the token"),
			Error:   err.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully generated token",
		Data: map[string]interface{}{
			"token": token.Hash,
		},
	})
}
