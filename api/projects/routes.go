package projects

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {
	group := sg.Group("/projects")
	group.POST("/validate-input", ValidateInputHandler)
}
