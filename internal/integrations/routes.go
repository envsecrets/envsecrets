package integrations

import (
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/integrations/:" + commons.INTEGRATION_TYPE)

	group.POST("/callback", nil)
}
