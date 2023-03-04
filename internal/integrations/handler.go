package integrations

import (
	"errors"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/labstack/echo/v4"
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
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Run the service handler.
	if err := service.Setup(ctx, serviceType, c.QueryParams()); err != nil {
		return c.JSON(http.StatusBadGateway, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to execute service callback",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully executed service callback",
	})
}

func ListEntitiesHandler(c echo.Context) error {

	//	Extract the entity type
	integration_type := c.Param(commons.INTEGRATION_TYPE)
	serviceType := commons.IntegrationType(integration_type)
	if !serviceType.IsValid() {
		return errors.New("invalid integration type")
	}

	/* 	//	Validate required query params.
	   	requiredParams := []string{"org_id"}
	   	for _, item := range requiredParams {
	   		if c.QueryParam(item) == "" {
	   			return errors.New("invalid or incomplete query parameters")
	   		}
	   	}
	*/
	//	Get the service.
	service := GetService()

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

	//	Run the service handler.
	entities, err := service.ListEntities(ctx, client, serviceType, c.Param(commons.INTEGRATION_ID))
	if err != nil {
		return c.JSON(http.StatusBadGateway, &commons.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to execute service callback",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &commons.APIResponse{
		Code:    http.StatusOK,
		Message: "successfully executed service callback",
		Data:    entities,
	})
}
