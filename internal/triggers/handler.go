package triggers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	inviteCommons "github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/mail"
	"github.com/envsecrets/envsecrets/internal/mail/commons"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/internal/organisations"
	organisationCommons "github.com/envsecrets/envsecrets/internal/organisations/commons"
	permissionCommons "github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/roles"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
	"github.com/labstack/echo/v4"
)

const (
	KEY_BYTES = 32
)

/*
// Called when a new row is inserted inside the `secrets` table.
func SecretInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row secretCommons.Secret
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Since the newly inserted secret is already base64 encoded,
	//	mark it 'encoded'
	row.MarkEncoded()

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	--- Flow ---
	//	1. Get the events linked to this new secret row.
	//	2. Call the appropriate integration service to sync the secrets.
	events, err := events.GetBySecret(ctx, client, row.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get events associated with this secret",
			Error:   err.Error(),
		})
	}

	if len(*events) == 0 {
		return c.JSON(http.StatusOK, &clients.APIResponse{
			Message: "there are no events in this environment to sync this secret with",
		})
	}

	//	Get the organisation to which these secrets belong to.
	organisation, err := organisations.GetService().GetByEnvironment(ctx, client, row.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get organisation to which these secrets are associated",
			Error:   err.Error(),
		})
	}

	//	Decrypt the value of every secret.
	decrypted, err := secrets.Decrypt(ctx, client, &secretCommons.DecryptOptions{
		OrgID:  organisation.ID,
		Secret: &row,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the secrets",
			Error:   err.Error(),
		})
	}

	//	Get the integration service
	integrationService := integrations.GetService()

	for _, event := range *events {
		if err := integrationService.Sync(ctx, client, &integrationCommons.SyncOptions{
			IntegrationID: event.Integration.ID,
			EventID:       event.ID,
			EntityDetails: event.EntityDetails,
			Data:          &decrypted.Data,
		}); err != nil {
			log.Printf("failed to push secret with ID %s for %s event: %s", row.ID, event.Integration.Type, event.ID)
			log.Println(err)
		}
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully synced secrets",
	})
}

// Called when a new row is inserted inside the `events` table.
func EventInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row eventCommons.Event
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	--- Flow ---
	//	1. Get the organisation linked to the environment of this event.
	//	2. Get the integration linked to this event.
	//	3. Fetch latest secrets linked to the environment of this event.
	//	4. Call the appropriate integration service to sync the secrets.

	//	Get the organisation to which this event's environment belong to.
	organisation, err := organisations.GetService().GetByEnvironment(ctx, client, row.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get organisation to which this event are associated",
			Error:   err.Error(),
		})
	}

	response, err := secrets.Get(ctx, client, &secretCommons.GetOptions{
		EnvID: row.EnvID,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get secrets associated with this event",
			Error:   err.Error(),
		})
	}

	//	Decrypt the value of every secret.
	decrypted, err := secrets.Decrypt(ctx, client, &secretCommons.DecryptOptions{
		OrgID:  organisation.ID,
		Secret: response,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to decrypt the secrets",
			Error:   err.Error(),
		})
	}

	//	Get the integration service
	integrationService := integrations.GetService()

	if err := integrationService.Sync(ctx, client, &integrationCommons.SyncOptions{
		IntegrationID: row.IntegrationID,
		EventID:       row.ID,
		EntityDetails: row.EntityDetails,
		Data:          &decrypted.Data,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: fmt.Sprintf("Failed to push secret with ID %s for event: %s", row.ID, row.ID),
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully synced secrets",
	})
}

*/ // Called when a new row is inserted inside the `secrets` table.
func SecretDeleteLegacy(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row secretCommons.Secret
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Get the organisation ID.
	organisation, err := organisations.GetService().GetByEnvironment(ctx, client, row.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get the organisation to which this environment is associated",
			Error:   err.Error(),
		})
	}

	//	Get subscriptions for this organisation
	orgSubscriptions, err := subscriptions.List(ctx, client, &subscriptions.ListOptions{OrgID: organisation.ID})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to get the subscriptions for this organisation",
			Error:   err.Error(),
		})
	}

	var active bool
	for _, item := range *orgSubscriptions {
		if item.Status == subscriptions.StatusActive {
			active = true
			break
		}
	}

	//	If no subscription is active,
	//	only keep the latest 5 version active.
	if !active {

		cleanupUntilVersion := *row.Version - 5
		if cleanupUntilVersion > 0 {
			if err := secrets.Cleanup(ctx, client, &secretCommons.CleanupSecretOptions{
				EnvID:   row.EnvID,
				Version: cleanupUntilVersion,
			}); err != nil {
				return c.JSON(http.StatusBadRequest, &clients.APIResponse{
					Message: "Failed to delete older versions of this secret",
					Error:   err.Error(),
				})
			}
		}
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully deleted legacy secrets",
	})
}

// Called when a new row is inserted inside the `users` table.
func UserInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var user userCommons.User
	if err := MapToStruct(payload.Event.Data.New, &user); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Shoot a welcome email to the user
	if err := mail.GetService().SendWelcomeEmail(ctx, &user); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to send welcome email",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "trigger completed successfully",
	})
}

// Called when a new row is inserted inside the `organisations` table.
func OrganisationCreated(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
			Error:   err.Error(),
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row organisationCommons.Organisation
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Generate default roles for the organisation.
	_, err := roles.Insert(ctx, client, &roles.RoleInsertOptions{
		OrgID: row.ID,
		Name:  "viewer",
		Permissions: permissionCommons.Permissions{
			Projects: permissionCommons.CRUD{
				Read: true,
			},
		},
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: fmt.Sprintf("Failed to create role: %s", "viewer"),
			Error:   err.Error(),
		})
	}

	_, err = roles.Insert(ctx, client, &roles.RoleInsertOptions{
		OrgID: row.ID,
		Name:  "editor",
		Permissions: permissionCommons.Permissions{
			Projects: permissionCommons.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Environments: permissionCommons.CRUD{
				Create: true,
				Update: true,
				Delete: true,
			},
		},
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: fmt.Sprintf("Failed to create role: %s", "editor"),
			Error:   err.Error(),
		})
	}

	adminRole, err := roles.Insert(ctx, client, &roles.RoleInsertOptions{
		OrgID: row.ID,
		Name:  "admin",
		Permissions: permissionCommons.Permissions{
			Integrations: permissionCommons.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Permissions: permissionCommons.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Projects: permissionCommons.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Environments: permissionCommons.CRUD{
				Create: true,
				Update: true,
				Delete: true,
			},
		},
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: fmt.Sprintf("Failed to create role: %s", "admin"),
			Error:   err.Error(),
		})
	}

	//	Generate a symmetric key for cryptographic operations in this organisation.
	keyBytes, err := globalCommons.GenerateRandomBytes(KEY_BYTES)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
			Message: "Failed to generate symmetric key for this org",
			Error:   err.Error(),
		})
	}

	//	Encrypt the key using owner's public key
	publicKeyBytes, err := keys.GetPublicKeyByUserID(ctx, client, row.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
			Message: "Failed to fetch owner's public key",
			Error:   err.Error(),
		})
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyBytes)
	result, err := keys.SealAsymmetricallyAnonymous(keyBytes, publicKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to seal org's symmetric key with owner's public key",
			Error:   err.Error(),
		})
	}

	if err := memberships.CreateWithUserID(ctx, client, &memberships.CreateOptions{
		UserID: row.UserID,
		OrgID:  row.ID,
		RoleID: adminRole.ID,
		Key:    base64.StdEncoding.EncodeToString(result),
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "Failed to add the owner as an admin member",
			Error:   err.Error(),
		})
	}

	/* 	//	Add envsecrets as a bot
	   	botPublicKeyBytes, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PUBLIC_KEY"))
	   	if err != nil {
	   		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
	   			Message: "Failed to base64 decode server's public key",
	   			Error:   err.Error(),
	   		})
	   	}

	   	var botPublicKey [32]byte
	   	copy(botPublicKey[:], botPublicKeyBytes)
	   	result, err = keys.SealAsymmetricallyAnonymous(keyBytes, botPublicKey)
	   	if err != nil {
	   		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
	   			Message: "Failed to seal org's symmetric key with bot's public key",
	   			Error:   err.Error(),
	   		})
	   	}

	   	if err := organisations.UpdateServerKeyCopy(ctx, client, &organisations.UpdateServerKeyCopyOptions{
	   		OrgID: row.ID,
	   		Key:   base64.StdEncoding.EncodeToString(result),
	   	}); err != nil {
	   		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
	   			Message: "Failed to save server's key copy",
	   			Error:   err.Error(),
	   		})
	   	}
	*/
	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully generated symmetric key and created default roles",
	})
}

// Called when a new row is inserted inside the `projects` table.
func ProjectInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row projects.Project
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Get the environments service.
	service := environments.GetService()

	//	Create default environments for this new project.
	envs := []string{"dev", "test", "staging", "prod"}
	for _, item := range envs {
		if _, err := service.CreateWithUserID(ctx, client, &environments.CreateOptions{
			Name:      item,
			ProjectID: row.ID,
			UserID:    row.UserID,
		}); err != nil {
			return c.JSON(http.StatusBadRequest, &clients.APIResponse{
				Message: "Failed to create the environment: " + item,
				Error:   err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully created default environments",
	})
}

// Called when a new row is inserted inside the `invites` table.
func InviteInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row inviteCommons.Invite
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext, EchoContext: c})

	//	Get the mailer service.
	service := mail.GetService()

	//	Send the invitation email.
	if err := service.Invite(ctx, &commons.InvitationOptions{
		ID:            row.ID,
		ReceiverEmail: row.Email,
		OrgID:         row.OrgID,
		SenderID:      row.UserID,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, &clients.APIResponse{
			Message: "Failed to send the invitation email",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &clients.APIResponse{
		Message: "successfully sent invitation email to " + row.Email,
	})
}
