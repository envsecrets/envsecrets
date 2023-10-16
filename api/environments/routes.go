package environments

import (
	"github.com/labstack/echo/v4"
)

const (
	ENV_ID = "env_id"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/environments")
	group.POST("/validate-input", ValidateInputHandler)

	environment := group.Group("/:" + ENV_ID)
	environment.POST("/sync-password", SyncWithPasswordHandler)
	environment.POST("/sync", SyncHandler)
}
