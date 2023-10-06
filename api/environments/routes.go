package environments

import (
	"github.com/labstack/echo/v4"
)

const (
	ENV_ID = "env_id"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/environments/:" + ENV_ID)

	group.POST("/sync-password", SyncWithPasswordHandler)
	group.POST("/sync", SyncHandler)
}
