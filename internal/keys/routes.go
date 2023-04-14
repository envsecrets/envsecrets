package keys

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/keys")

	group.GET("/backup", KeyBackupHandler)
	group.POST("/restore", KeyRestoreHandler)
}
