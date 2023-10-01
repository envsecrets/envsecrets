package invites

import (
	"net/http"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func AcceptHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.AcceptRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Call the service function.
	if err := GetService().Accept(ctx, client, payload.ID); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to accept the invite",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully accepted the invite",
	})
}

// ---	Flow ---
// 1. Fetch the organisation's symmetric key using the organisation ID.
// 2. Decrypt the organisation's symmetric key using user's received password.
// 3. Pull the public key of the invitee.
// 4. Create a new copy of the organisation's symmetric key encrypted with the invitee's public key.
// 5. Save the new encrypted copy in the database and send the invite.
func SendHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload commons.SendRequestOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Extract the user's email from JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.Claims)

	//	Decrypt and get the bytes of user's own copy of organisation's encryption key.
	key, err := keys.DecryptMemberKey(ctx, client, claims.Hasura.UserID, &keyCommons.DecryptOptions{
		OrgID:    payload.OrgID,
		Password: payload.InviterPassword,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the secrets",
			Error:   err.Error(),
		})
	}

	//	Call the service function.
	if err := GetService().Send(ctx, client, &commons.SendOptions{
		OrgID:        payload.OrgID,
		RoleID:       payload.RoleID,
		InviteeEmail: payload.InviteeEmail,
		InviterID:    claims.Hasura.UserID,
		Key:          key,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to send the invite",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully sent the invite",
	})
}
