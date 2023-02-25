package events

import (
	"github.com/envsecrets/envsecrets/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	events := sg.Group("/events")

	//	envsecrets permissions group
	permissions := events.Group("/permissions")

	//	Prepend the webhook middleware to this group.
	permissions.Use(middlewares.HasWebhookHeader)

	permissions.POST("/organisation", OrganisationLevelPermissions)
	permissions.POST("/organisation/new", OrganisationInserted)
	permissions.POST("/organisation/delete", OrganisationDeleted)
	permissions.POST("/project", ProjectLevelPermissions)
	permissions.POST("/environment", EnvironmentLevelPermissions)
}
