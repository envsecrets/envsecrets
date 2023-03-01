package triggers

import (
	"github.com/envsecrets/envsecrets/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	triggers := sg.Group("/triggers")

	//	Prepend the webhook middleware to this group.
	triggers.Use(middlewares.WebhookHeader())

	//	envsecrets secrets group
	secrets := triggers.Group("/secrets")

	secrets.POST("/new", SecretInserted)

	//	envsecrets permissions group
	permissions := triggers.Group("/permissions")

	permissions.POST("/organisation", OrganisationLevelPermissions)
	permissions.POST("/organisation/new", OrganisationInserted)
	permissions.POST("/organisation/delete", OrganisationDeleted)
	permissions.POST("/project", ProjectLevelPermissions)
	permissions.POST("/environment", EnvironmentLevelPermissions)
}
