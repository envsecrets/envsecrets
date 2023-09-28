package environments

import (
	"github.com/envsecrets/envsecrets/internal/environments/commons"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/environments/:" + commons.ENV_ID)

	group.POST("/sync-password", SyncWithPasswordHandler)
	group.POST("/sync", SyncHandler)
}
