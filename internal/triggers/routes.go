package triggers

import (
	"github.com/envsecrets/envsecrets/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func AddRoutes(sg *echo.Group) {

	triggers := sg.Group("/triggers")

	//	Prepend the webhook middleware to this group.
	triggers.Use(middlewares.WebhookHeader())

	//	secrets group
	secrets := triggers.Group("/secrets")

	secrets.POST("/new", SecretInserted)
	secrets.POST("/delete-legacy", SecretDeleteLegacy)

	//	events group
	events := triggers.Group("/events")

	events.POST("/new", EventInserted)

	//	users group
	//users := triggers.Group("/users")
	//users.POST("/new", UserInserted)

	//	invites group
	invites := triggers.Group("/invites")
	invites.POST("/new", InviteInserted)

	//	organisations group
	organisations := triggers.Group("/organisations")

	organisations.POST("/new", OrganisationCreated)

	//	projects group
	projects := triggers.Group("/projects")
	projects.POST("/new", ProjectInserted)
}
