package invites

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/internal/organisations"
	organisationCommons "github.com/envsecrets/envsecrets/internal/organisations/commons"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/golang-jwt/jwt/v4"
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
		return c.String(http.StatusNotFound, "Invite not found")
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

	//	Get the server's copy of org-key.
	serverOrgKey, err := organisations.GetService().GetServerKeyCopy(ctx, client, invite.OrgID)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Failed to fetch server's copy of org-key")
	}

	//	Decrypt the copy with server's private key (in env vars).
	serverPrivateKey, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PRIVATE_KEY"))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Failed to base64 decode server's private key.")
	}
	serverPublicKey, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PUBLIC_KEY"))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Failed to base64 decode server's private key.")
	}

	result, err := keys.DecryptAsymmetricallyAnonymous(serverPublicKey, serverPrivateKey, serverOrgKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
			Message: "Failed to decrypt server's copy of org-key",
			Error:   err.Error(),
		})
	}

	//	Fetch the invitee's public key.
	inviteePublicKeyBytes, err := keys.GetPublicKeyByUserID(ctx, client, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch invitee's public key",
			Error:   err.Error(),
		})
	}
	var invitePublicKey [32]byte
	copy(invitePublicKey[:], inviteePublicKeyBytes)

	//	Create key copy for the invitee.
	result, err = keys.SealAsymmetricallyAnonymous(result, invitePublicKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
			Message: "Failed to create a copy of org-key for the invitee",
			Error:   err.Error(),
		})
	}

	if err := memberships.CreateWithUserID(ctx, client, &memberships.CreateOptions{
		UserID: user.ID,
		OrgID:  invite.OrgID,
		RoleID: invite.RoleID,
		Key:    base64.StdEncoding.EncodeToString(result),
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to add the invitee as a member",
			Error:   err.Error(),
		})
	}

	//	Mark the invite "accepted".
	if _, err := service.Update(ctx, client, id, &commons.UpdateOptions{
		Accepted: true,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to accept the invite",
			Error:   err.Error(),
		})
	}

	//	Reduce the invite limit in organisation by 1.
	if err := organisations.GetService().UpdateInviteLimit(ctx, client, &organisationCommons.UpdateInviteLimitOptions{
		ID:               invite.OrgID,
		IncrementLimitBy: -1,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrement org's invite limit",
			Error:   err.Error(),
		})
	}

	//	Redirect the user to homepage.
	return c.Redirect(http.StatusPermanentRedirect, os.Getenv("FE_URL")+"/choose-organisation")
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
