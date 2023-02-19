package integrations

import (
	"errors"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/labstack/echo/v4"
)

func CallbackHandler(c echo.Context) error {

	//	Extract the entity type
	integration_type := c.Param(commons.INTEGRATION_TYPE)
	serviceType := commons.IntegrationType(integration_type)
	if !serviceType.IsValid() {
		return errors.New("invalid integration type")
	}

	//	Get the service.
	service := GetService()

	//	Run the service handler.
	if err := service.Callback(serviceType, c.QueryParams()); err != nil {
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
