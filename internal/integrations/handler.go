package integrations

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/labstack/echo/v4"
)

func SetupCallbackHandler(c echo.Context) error {

	//	Extract the entity type
	integration_type := c.Param(commons.INTEGRATION_TYPE)
	serviceType := commons.IntegrationType(integration_type)
	if !serviceType.IsValid() {
		return errors.New("invalid integration type")
	}

	//	Get the service.
	service := GetService()

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	params := c.QueryParams()
	state := params.Get("state")

	//	Extract the Organisation ID and Authorization token from State.
	payload := strings.Split(state, "/")
	if len(payload) != 2 {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "invalid callback state",
			Error:   "invalid callback state",
		})
	}
	orgID := payload[0]
	token := payload[1]

	//	Initialize Hasura client with token extract from state parameter
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   string(clients.AuthorizationHeader),
				Value: "Bearer " + token,
			},
		},
	})

	options := make(map[string]interface{})

	for key, value := range params {
		options[key] = value[0]
	}

	//	Run the service handler.
	_, err := service.Setup(ctx, client, serviceType, &commons.SetupOptions{
		OrgID:   orgID,
		Options: options,
	})
	if err != nil {
		return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/integrations/catalog?setup_action=install&setup_status=failed&integration_type=%s", os.Getenv("FE_URL"), integration_type))
	}

	//	Redirect the user to front-end to complete post-integration steps.
	return c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s/integrations/%s/sync?setup_action=install&setup_status=successful", os.Getenv("FE_URL"), integration_type))
}

func SetupHandler(c echo.Context) error {

	//	Extract the entity type
	integration_type := c.Param(commons.INTEGRATION_TYPE)
	serviceType := commons.IntegrationType(integration_type)
	if !serviceType.IsValid() {
		return errors.New("invalid integration type")
	}

	//	Unmarshal the incoming payload
	var payload commons.SetupOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Get the service.
	service := GetService()

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Run the service handler.
	_, err := service.Setup(ctx, client, serviceType, &payload)
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Message: err.GenerateMessage("failed to connect integration"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully integrated " + integration_type,
	})
}

func ListEntitiesHandler(c echo.Context) error {

	//	Extract the entity type
	integration_type := c.Param(commons.INTEGRATION_TYPE)
	serviceType := commons.IntegrationType(integration_type)
	if !serviceType.IsValid() {
		return errors.New("invalid integration type")
	}

	//	Get the service.
	service := GetService()

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	options := make(map[string]interface{})

	for key, value := range c.QueryParams() {
		options[key] = value[0]
	}

	//	Run the service handler.
	entities, err := service.ListEntities(ctx, client, serviceType, c.Param(commons.INTEGRATION_ID), options)
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Message: err.GenerateMessage("failed to fetch integration entities"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully fetched integration entities",
		Data:    entities,
	})
}

func ListSubEntitiesHandler(c echo.Context) error {

	//	Extract the entity type
	integration_type := c.Param(commons.INTEGRATION_TYPE)
	serviceType := commons.IntegrationType(integration_type)
	if !serviceType.IsValid() {
		return errors.New("invalid integration type")
	}

	//	Get the service.
	service := GetService()

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Run the service handler.
	entities, err := service.ListSubEntities(ctx, client, serviceType, c.Param(commons.INTEGRATION_ID), c.QueryParams())
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &clients.APIResponse{
			Message: err.GenerateMessage("failed to fetch integration entities"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully fetched integration entities",
		Data:    entities,
	})
}
