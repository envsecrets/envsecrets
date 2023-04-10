package payments

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	payments := sg.Group("/payments")

	//	Client group.
	//	To be called using `Authorization` header / JWT.
	client := payments.Group("/client")

	//	Set the group's routes.
	client.GET("/session", CreateCheckoutSession)

	//	Server group.
	//	To be called with webhook header.
	server := payments.Group("/server")

	//	Set the group's routes.
	server.POST("/webhook", CheckoutWebhook)
}
