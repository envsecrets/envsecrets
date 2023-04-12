package tokens

import (
	"net/http"
	"time"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations"
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

	//	Initialize Hasura client with admin privileges
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

	//	Get the organisation for this environment.
	organisation, err := organisations.GetByEnvironment(ctx, client, payload.EnvID)
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to get the organisation this environment is associated with"),
			Error:   err.Error.Error(),
		})
	}

	token, err := Create(ctx, client, &commons.CreateServiceOptions{
		OrgID:  organisation.ID,
		EnvID:  payload.EnvID,
		Expiry: expiry,
	})
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
			"token": token,
		},
	})
}
