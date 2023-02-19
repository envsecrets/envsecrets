package secrets

import "github.com/labstack/echo/v4"

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/secrets")

	group.POST("/", SetHandler)
	group.GET("/", GetHandler)
}
