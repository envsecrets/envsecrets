package tokens

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	commonGroup := sg.Group("/tokens")

	commonGroup.GET("/accept", nil)
}
