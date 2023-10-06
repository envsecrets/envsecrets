package invites

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/invites")

	group.POST("", SendHandler)
	group.POST("/accept", AcceptHandler)
}
