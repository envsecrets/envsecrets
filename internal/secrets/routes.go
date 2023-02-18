package secrets

import "github.com/labstack/echo/v4"

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/secrets")

	group.POST("/keys/generate", nil)
	group.POST("/set", nil)
	group.GET("/get", nil)
}
