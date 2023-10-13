package events

import (
	"net/http"

	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/events"
	"github.com/labstack/echo/v4"
)

func ActionsGetHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload clients.HasuraActionRequestPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Extract the arguments out of action payload.
	var inputs events.ActionsGetOptions
	if err := globalCommons.MapToStruct(payload.Input.Args, &inputs); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the inputs",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	events, err := events.GetByEnvironment(ctx, client, inputs.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to get the events",
			Error:   err.Error(),
		})
	}

	//	Transform the events to the response format
	var response []map[string]interface{}

	for _, event := range *events {
		payload := make(map[string]interface{})
		payload["id"] = event.ID
		payload["title"] = event.Integration.GetTitle()
		payload["type"] = event.Integration.Type
		//payload["description"] = event.Integration.GetDescription()
		//payload["subtitle"] = event.Integration.GetSubtitle()
		payload["link"] = event.GetEntityLink()
		payload["name"] = event.GetEntityTitle()
		payload["entity_type"] = event.GetEntityType()
		response = append(response, payload)
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Data: response,
	})
}
