package integrations

import (
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	commonGroup := sg.Group("/integrations/:" + commons.INTEGRATION_TYPE)

	commonGroup.GET("/setup", SetupHandler)

	integrationsGroup := commonGroup.Group("/:" + commons.INTEGRATION_ID)
	integrationsGroup.GET("/entities", ListEntitiesHandler)
	integrationsGroup.POST("/trigger", nil, middlewares.WebhookHeader())
}
