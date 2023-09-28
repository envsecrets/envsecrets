package events

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/events")
	actions := group.Group("/actions")

	actions.POST("/get", ActionsGetHandler)
}
