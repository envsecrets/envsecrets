package tokens

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	keysCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/tokens"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload CreateOptions
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with user's token
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetService().GetByEnvironment(ctx, client, payload.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the organisation this environment is associated with",
			Error:   err.Error(),
		})
	}

	//	Initialize the tokens service.
	service := tokens.GetService()

	//	Validate abuse limits
	if err := validateInput(ctx, client, payload.EnvID); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Your current plan does not allow creating more tokens.",
			Error:   err.Error(),
		})
	}

	//	Extract the user's email from JWT
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.Claims)

	//	Decrypt and get the bytes of user's own copy of organisation's encryption key.
	orgKey, err := keys.DecryptMemberKey(ctx, client, claims.Hasura.UserID, &keysCommons.DecryptOptions{
		OrgID:    organisation.ID,
		Password: payload.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the organisation's encryption key. Maybe, entered password is invalid.",
			Error:   err.Error(),
		})
	}

	//	Create the token
	expiry, err := time.ParseDuration(payload.Expiry)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to parse expiry duration",
			Error:   err.Error(),
		})
	}

	token, err := service.Create(ctx, client, &tokens.CreateOptions{
		OrgKey: orgKey,
		EnvID:  payload.EnvID,
		Expiry: expiry,
		Name:   payload.Name,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to create the token",
			Error:   err.Error(),
		})
	}

	//	Encode the token
	hash := hex.EncodeToString(token)

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully generated token",
		Data: map[string]interface{}{
			"token": hash,
		},
	})
}
