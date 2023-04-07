package integrations

import (
	"errors"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func SetupHandler(c echo.Context) error {

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

	//	Run the service handler.
	if err := service.Setup(ctx, serviceType, c.QueryParams()); err != nil {
		log.Error(err)
		return c.Redirect(http.StatusPermanentRedirect, os.Getenv("FE_URL")+"/integrations?setup_action=install&setup_status=failed&integration_type="+integration_type)
	}

	//	Redirect the user to front-end to complete post-integration steps.
	return c.Redirect(http.StatusPermanentRedirect, os.Getenv("FE_URL")+"/integrations?setup_action=install&setup_status=successful&integration_type="+integration_type)
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

	//	Run the service handler.
	entities, err := service.ListEntities(ctx, client, serviceType, c.Param(commons.INTEGRATION_ID))
	if err != nil {
		return c.JSON(err.Type.GetStatusCode(), &commons.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("failed to fetch integration entities"),
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully fetched integration entities",
		Data:    entities,
	})
}
