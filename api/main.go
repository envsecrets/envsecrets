package api

import (
	"github.com/envsecrets/envsecrets/api/auth"
	"github.com/envsecrets/envsecrets/api/environments"
	"github.com/envsecrets/envsecrets/api/events"
	"github.com/envsecrets/envsecrets/api/integrations"
	"github.com/envsecrets/envsecrets/api/invites"
	"github.com/envsecrets/envsecrets/api/payments"
	"github.com/envsecrets/envsecrets/api/tokens"
	"github.com/envsecrets/envsecrets/api/triggers"
	"github.com/envsecrets/envsecrets/internal/actions"
	"github.com/envsecrets/envsecrets/internal/secrets"
	"github.com/labstack/echo/v4"
)

// AddRoutes adds all routes to the echo instance.
func AddRoutes(e *echo.Echo) {

	//	API	Version 1 Group
	v1Group := e.Group("/v1")

	triggers.AddRoutes(v1Group)
	actions.AddRoutes(v1Group)
	auth.AddRoutes(v1Group)
	secrets.AddRoutes(v1Group)
	integrations.AddRoutes(v1Group)
	environments.AddRoutes(v1Group)
	invites.AddRoutes(v1Group)
	payments.AddRoutes(v1Group)
	tokens.AddRoutes(v1Group)
	events.AddRoutes(v1Group)
	//keys.AddRoutes(v1Group)
}
