package secrets

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	group := sg.Group("/secrets")
	group.Use(middlewares.JWTAuth(func(c echo.Context) bool {
		return c.Request().Method == echo.GET
	}))

	group.POST("", SetHandler)
	group.DELETE("", DeleteHandler)

	//	Custom middlewares for GET routes
	middlewares := []echo.MiddlewareFunc{
		middlewares.TokenHeader(),
		middlewares.JWTAuth(func(c echo.Context) bool {

			//	Skip the middleware if it has an environment token.
			return c.Request().Header.Get(string(clients.TokenHeader)) != ""
		}),
	}

	group.GET("", GetHandler, middlewares...)
	//	group.GET("/values", GetHandler, middlewares...)
}
