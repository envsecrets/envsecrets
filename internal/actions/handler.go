package actions

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/labstack/echo/v4"
)

func EnvironmentCreate(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraActionPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Unmarshal the data interface to our required entity.
	var options environments.CreateOptions
	if err := MapToStruct(payload.Input.Args, &options); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	--- Checks to perform ---
	//
	//	1. Organisation is active on at least 1 paid plan.

	//	Fetch the org_id by using project_id.
	project, err := projects.Get(ctx, client, options.ProjectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the project details",
			Error:   err.Error(),
		})
	}

	activeSubscriptions, err := subscriptions.List(ctx, client, &subscriptions.ListOptions{OrgID: project.OrgID})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the subscription details",
			Error:   err.Error(),
		})
	}

	var allowed bool
	for _, item := range *activeSubscriptions {
		if item.Status == subscriptions.StatusActive {
			allowed = true
			break
		}
	}

	if !allowed {
		return c.JSON(http.StatusForbidden, &HasuraActionErrorResponse{
			Message: "Your current plan does not allow custom environments. Please upgrade your plan!",
		})
	}

	//	Create the environment
	environment, err := environments.Create(ctx, client, &options)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to create the environment",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, environment)
}
