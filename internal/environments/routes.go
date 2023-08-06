package environments

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/environments/:" + ENV_ID)

	group.POST("/sync", SyncHandler)
}
