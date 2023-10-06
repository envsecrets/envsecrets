package integrations

import (
	"github.com/envsecrets/envsecrets/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	commonGroup := sg.Group("/integrations/:" + INTEGRATION_TYPE)

	commonGroup.GET("/callback/setup", SetupCallbackHandler)
	commonGroup.POST("/setup", SetupHandler)

	integrationsGroup := commonGroup.Group("/:" + INTEGRATION_ID)
	integrationsGroup.GET("/entities", ListEntitiesHandler)
	integrationsGroup.GET("/sub-entities", ListSubEntitiesHandler)
	integrationsGroup.POST("/trigger", nil, middlewares.WebhookHeader())
}
