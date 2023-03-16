package secrets

import "github.com/labstack/echo/v4"

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/secrets")

	group.POST("", SetHandler)
	group.POST("/merge", MergeHandler)
	group.GET("", GetHandler)
	group.DELETE("", DeleteHandler)

	keys := group.Group("/keys")
	keys.GET("/backup", KeyBackupHandler)
	keys.POST("/restore", KeyRestoreHandler)
}
