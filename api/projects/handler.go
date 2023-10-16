package projects

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/envsecrets/envsecrets/utils"
	"github.com/labstack/echo/v4"
)

func ValidateInputHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload clients.HasuraInputValidationPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
			Message: "failed to parse the body",
			Extensions: &clients.HasuraActionsResponseExtensions{
				Error: err,
			},
		})
	}

	//	Unmarshal the data interface to our required entity.
	var rows []projects.Project
	if err := utils.MapToStruct(payload.Data.Input, &rows); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
			Message: "failed to unmarshal new data",
			Extensions: &clients.HasuraActionsResponseExtensions{
				Error: err,
			},
		})
	}

	// Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	// Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Check the number of existing projects for the organisation.
	//	If the number of projects is greater than the allowed limit, proceed to check whether the organisation has an active subscription.
	//	Otherwise, approve the inputs and allow for creation of the project.
	for _, row := range rows {
		projects, err := projects.GetService().List(ctx, client, &projects.ListOptions{
			OrgID: row.OrgID,
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: "failed to get the projects",
				Extensions: &clients.HasuraActionsResponseExtensions{
					Error: err,
				},
			})
		}

		//	If the number of projects is greater than the allowed limit, proceed to check whether the organisation has an active subscription.
		//	Otherwise, approve the inputs and allow for creation of the project.
		if len(projects) < FREE_TIER_LIMIT_NUMBER_OF_PROJECTS {
			continue
		}

		//	Validate whether the organisation an active premium subscription.
		//	We do this by fetching the subscriptions by the organisation ID.
		//	We then check if any subscription is active.
		subscriptions, err := subscriptions.GetService().GetByOrgID(ctx, client, row.OrgID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: "failed to get the subscriptions",
				Extensions: &clients.HasuraActionsResponseExtensions{
					Error: err,
				},
			})
		}

		//	If there are no subscriptions, or if even a single subscription is not active, return an error.
		if len(*subscriptions) == 0 {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: clients.ErrBreachingAbuseLimit.Error(),
			})
		}

		active := subscriptions.IsActiveAny()
		if !active {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: clients.ErrBreachingAbuseLimit.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, &clients.HasuraActionResponse{
		Message: "inputs validated and permitted",
	})
}
