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
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/envsecrets/envsecrets/utils"
	"github.com/golang-jwt/jwt/v4"
	echo "github.com/labstack/echo/v4"
)

func ValidateInputHandler(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload clients.HasuraInputValidationPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
			Message: "failed to parse the body",
			Extensions: &clients.HasuraActionsResponseExtensions{
				Error: err,
			},
		})
	}

	//	Unmarshal the data interface to our required entity.
	var rows []environments.Environment
	if err := utils.MapToStruct(payload.Data.Input, &rows); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
			Message: "failed to unmarshal new data",
			Extensions: &clients.HasuraActionsResponseExtensions{
				Error: err,
			},
		})
	}

	// Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	// Initialize new Hasura client
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type:          clients.HasuraClientType,
		Authorization: c.Request().Header.Get(echo.HeaderAuthorization),
	})

	//	Check the number of existing environments for the organisation.
	//	If the number of environments is greater than the allowed limit, proceed to check whether the organisation has an active subscription.
	//	Otherwise, approve the inputs and allow for creation of the project.
	for _, row := range rows {
		environments, err := environments.GetService().List(ctx, client, &environments.ListOptions{
			ProjectID: row.ProjectID,
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: "failed to get the environments",
				Extensions: &clients.HasuraActionsResponseExtensions{
					Error: err,
				},
			})
		}

		//	If the number of environments is greater than the allowed limit, proceed to check whether the organisation has an active subscription.
		//	Otherwise, approve the inputs and allow for creation of the project.
		if len(environments) < FREE_TIER_LIMIT_NUMBER_OF_ENVIRONMENTS {
			continue
		}

		organisation, err := projects.GetService().GetOrganisation(ctx, client, row.ProjectID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: "failed to get the organisation",
			})
		}

		//	Validate whether the organisation an active premium subscription.
		//	We do this by fetching the subscriptions by the organisation ID.
		//	We then check if any subscription is active.
		subscriptions, err := subscriptions.GetService().GetByOrgID(ctx, client, organisation.ID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: "failed to get the subscriptions",
				Extensions: &clients.HasuraActionsResponseExtensions{
					Error: err,
				},
			})
		}

		//	If there are no subscriptions, or if even a single subscription is not active, return an error.
		if len(*subscriptions) == 0 {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: clients.ErrBreachingAbuseLimit.Error(),
			})
		}

		active := subscriptions.IsActiveAny()
		if !active {
			return c.JSON(http.StatusBadRequest, &clients.HasuraActionResponse{
				Message: clients.ErrBreachingAbuseLimit.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, &clients.HasuraActionResponse{
		Message: "inputs validated and permitted",
	})
}

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
	var payload SyncWithPasswordOptions
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
	claims := token.Claims.(*clients.Claims)

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

	//	Get the environments service.
	service := environments.GetService()

	//	Call the service function.
	if err := service.Sync(ctx, client, &environments.SyncOptions{
		EnvID:    envID,
		EventIDs: payload.EventIDs,
		Pairs:    &decrypted.Data,
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
	var payload SyncOptions
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
	if err := payload.Pairs.Decode(); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to decode the secrets",
			Error:   err.Error(),
		})
	}

	//	Extract the user's email from JWT
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.Claims)

	//	Get the user's sync key and decrypt it with server's own encryption key.
	syncKey, err := keys.GetSyncKeyByUserID(ctx, client, claims.Hasura.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to get the user's sync key",
			Error:   err.Error(),
		})
	}

	decryptedSyncKeyBytes, err := keys.OpenSymmetricallyByServer(syncKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to decrypt the user's sync key",
			Error:   err.Error(),
		})
	}

	//	Now decrypt the secrets using the decrypted sync key.
	var decryptedSyncKey [32]byte
	copy(decryptedSyncKey[:], decryptedSyncKeyBytes)
	if err := payload.Pairs.Decrypt(decryptedSyncKey); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to decrypt the secrets",
			Error:   err.Error(),
		})
	}

	//	Call the service function.
	if err := environments.GetService().Sync(ctx, client, &environments.SyncOptions{
		EnvID:    envID,
		Pairs:    payload.Pairs,
		EventIDs: payload.EventIDs,
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
