package events

import "github.com/labstack/echo/v4"

func AddRoutes(sg *echo.Group) {

	events := sg.Group("/events")

	//	envsecrets permissions group
	permissions := events.Group("/permissions")

	permissions.POST("/organisation", OrganisationLevelPermissions)
	permissions.POST("/organisation/new", OrganisationInserted)
	permissions.POST("/organisation/delete", OrganisationDeleted)
	permissions.POST("/project", ProjectLevelPermissions)
	permissions.POST("/environment", EnvironmentLevelPermissions)
}
