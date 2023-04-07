package invites

import (
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/permissions"
	permissionCommons "github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/labstack/echo/v4"
)

func AcceptHandler(c echo.Context) error {

	//	Extract the invite ID and Key
	id := c.QueryParam("id")
	key := c.QueryParam("key")

	//	Get the service.
	service := GetService()

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Fetch the invite
	invite, err := service.Get(ctx, client, id)
	if err != nil {
		return c.String(err.Type.GetStatusCode(), "Invite not found")
	}

	//	Validate the key for this invite.
	if key != invite.Key {
		return c.String(http.StatusUnauthorized, "Failed to authenticate this invite. Ask for a new one!")
	}

	//	Return error if the invite has already been accepted.
	if invite.Accepted {
		return c.String(http.StatusForbidden, "This invite has already been accepted. Ask for a new one!")
	}

	//	Fetch the user.
	user, err := users.GetByEmail(ctx, client, invite.Email)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Failed to fetch the user for whom this invitation is meant for. Create an envsecrets account and re-try accepting this invite.")
	}

	//	Insert the user with appropriate role in the organisation.
	if err := permissions.GetService().Insert(permissionCommons.OrgnisationLevelPermission, ctx, client, permissionCommons.OrganisationPermissionsInsertOptions{
		OrgID:  invite.OrgID,
		RoleID: invite.RoleID,
		UserID: user.ID,
	}); err != nil {
		return c.JSON(err.Type.GetStatusCode(), &commons.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to accept invite"),
			Error:   err.Error.Error(),
		})
	}

	//	Mark the invite "accepted".
	if _, err := service.Update(ctx, client, id, &commons.UpdateOptions{
		Accepted: true,
	}); err != nil {
		return c.JSON(err.Type.GetStatusCode(), &commons.APIResponse{
			Code:    err.Type.GetStatusCode(),
			Message: err.GenerateMessage("Failed to accept invite"),
			Error:   err.Error.Error(),
		})
	}

	//	Redirect the user to homepage.
	return c.Redirect(http.StatusPermanentRedirect, os.Getenv("FE_URL")+"/choose-organisation")
}
