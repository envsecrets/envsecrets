package invites

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	commonGroup := sg.Group("/invites")

	commonGroup.GET("/accept", AcceptHandler)
}
