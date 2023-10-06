package environments

import (
	"net/http"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/keys"
	keysCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/golang-jwt/jwt/v4"
	echo "github.com/labstack/echo/v4"
)

// --- Flow ---
//
//  1. Get the organisation ID linked to this environment.
//  2. Decrypt the organisation's encryption key.
//  3. Fetch the secrets of this environment.
//     - Fetch the latest version if no version is specified in request payload.
//  4. Decrypt the secrets using organisation's encryption key.
//  5. Fetch the events linked to this environment.
//  6. Call the integration service for each of the events to sync the secrets.
func SyncWithPasswordHandler(c echo.Context) error {

	//	Extract the entity type
	envID := c.Param(ENV_ID)
	if envID == "" {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "invalid environment ID",
			Error:   "invalid environment ID",
		})
	}

	//	Unmarshal the incoming payload
	var payload environments.SyncWithPasswordRequestOptions
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

	//	Fetch the organisation using environment ID.
	organisation, err := organisations.GetService().GetByEnvironment(ctx, client, envID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to fetch the organisation this environment is associated with",
			Error:   err.Error(),
		})
	}

	//	Extract the user's email from JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.Claims)

	//	Decrypt and get the bytes of user's own copy of organisation's encryption key.
	key, err := keys.DecryptMemberKey(ctx, client, claims.Hasura.UserID, &keysCommons.DecryptOptions{
		OrgID:    organisation.ID,
		Password: payload.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the secrets",
			Error:   err.Error(),
		})
	}

	//	Fetch the secrets.
	response, err := secrets.Get(ctx, client, &secretCommons.GetOptions{
		EnvID:   envID,
		Key:     payload.Key,
		Version: payload.Version,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get secrets associated with this event",
			Error:   err.Error(),
		})
	}

	//	Decrypt the value of every secret.
	decrypted, err := secrets.Decrypt(ctx, client, &secretCommons.DecryptOptions{
		Secret: response,
		Key:    key,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the secrets",
			Error:   err.Error(),
		})
	}

	//	Call the service function.
	if err := environments.Sync(ctx, client, &environments.SyncOptions{
		EnvID:           envID,
		IntegrationType: payload.IntegrationType,
		Secrets:         &decrypted.Data,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to sync the secrets",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully synced secrets",
	})
}

func SyncHandler(c echo.Context) error {

	//	Extract the entity type
	envID := c.Param(ENV_ID)
	if envID == "" {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "invalid environment ID",
			Error:   "invalid environment ID",
		})
	}

	//	Unmarshal the incoming payload
	var payload environments.SyncRequestOptions
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

	//	Decode the values before sending them further.
	if err := payload.Data.Decode(); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to decode the secrets",
			Error:   err.Error(),
		})
	}

	//	Call the service function.
	if err := environments.Sync(ctx, client, &environments.SyncOptions{
		EnvID:           envID,
		Secrets:         payload.Data,
		IntegrationType: payload.IntegrationType,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to sync the secrets",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully synced secrets",
	})
}
