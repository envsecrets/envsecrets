package organisations

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/organisations")
	group.POST("", CreateHandler)
}
