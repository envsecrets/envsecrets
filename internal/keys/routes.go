package keys

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/keys")

	group.GET("/public-key", nil)
}
