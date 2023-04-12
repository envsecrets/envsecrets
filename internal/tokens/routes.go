package tokens

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	tokens := sg.Group("/tokens")

	tokens.POST("", CreateHandler)
}
