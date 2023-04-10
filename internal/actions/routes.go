package actions

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	actions := sg.Group("/actions")

	//	environments group
	environments := actions.Group("/environments")

	environments.POST("", EnvironmentCreate)
}
