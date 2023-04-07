package triggers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/events"
	eventCommons "github.com/envsecrets/envsecrets/internal/events/commons"
	"github.com/envsecrets/envsecrets/internal/integrations"
	integrationCommons "github.com/envsecrets/envsecrets/internal/integrations/commons"
	inviteCommons "github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/mail"
	"github.com/envsecrets/envsecrets/internal/mail/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/permissions"
	permissionCommons "github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
	"github.com/labstack/echo/v4"
)

//	Called when a new row is inserted inside the `secrets` table.
func SecretInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row secretCommons.Secret
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Cleanup old secrets. Only keep latest 10 secrets.
	cleanupUntilVersion := row.Version - 10
	if cleanupUntilVersion > 10 {
		if err := secrets.Cleanup(ctx, client, &secretCommons.CleanupSecretOptions{
			EnvID:   row.EnvID,
			Version: cleanupUntilVersion,
		}); err != nil {
			log.Println("Failed to cleanup older secret rows: ", err)

			//	Don't exit.
		}
	}

	//	--- Flow ---
	//	1. Get the events linked to this new secret row.
	//	2. Call the appropriate integration service to sync the secrets.

	events, err := events.GetBySecret(ctx, client, row.ID)
	if err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to get events associated with this secret",
			Error:   err.Error.Error(),
		})
	}

	data := make(map[string]secretCommons.Payload)

	if len(*events) > 0 {

		//	Get the organisation to which these secrets belong to.
		organisation, err := organisations.GetByEnvironment(ctx, client, row.EnvID)
		if err != nil {
			return c.JSON(http.StatusBadGateway, &APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to get organisation to which these secrets are associated",
				Error:   err.Error.Error(),
			})
		}

		//	Decrypt the value of every secret.
		for key, payload := range row.Data {

			//	If the secret is of type `ciphertext`,
			//	we will need to decode it first.
			if payload.Type == secretCommons.Ciphertext {
				secret, err := secrets.Decrypt(ctx, &secretCommons.DecryptSecretOptions{
					Data: secretCommons.Data{
						Key:     key,
						Payload: payload,
					},
					KeyLocation: organisation.ID,
					EnvID:       row.EnvID,
				})
				if err != nil {
					return c.JSON(http.StatusBadGateway, &APIResponse{
						Code:    http.StatusBadRequest,
						Message: "failed to decrypt value of secret: " + key,
						Error:   err.Error.Error(),
					})
				}

				//	Base64 decode the secret value
				b64Decoded, er := base64.StdEncoding.DecodeString(secret.Data.Plaintext)
				if er != nil {
					return c.JSON(http.StatusBadGateway, &APIResponse{
						Code:    http.StatusBadRequest,
						Message: "failed to base 64 decode the decrypted value of secret: " + key,
						Error:   er.Error(),
					})
				}

				payload.Value = string(b64Decoded)

			} else if payload.Type == secretCommons.Plaintext {

				//	Base64 decode the secret value
				b64Decoded, er := base64.StdEncoding.DecodeString(payload.Value.(string))
				if er != nil {
					return c.JSON(http.StatusBadGateway, &APIResponse{
						Code:    http.StatusBadRequest,
						Message: "failed to base 64 decode the decrypted value of secret: " + key,
						Error:   er.Error(),
					})
				}

				payload.Value = string(b64Decoded)
			}

			data[key] = payload
		}
	}

	//	Get the integration service
	integrationService := integrations.GetService()

	var wg sync.WaitGroup
	for _, event := range *events {
		wg.Add(1)
		go func(event *eventCommons.Event) {
			if err := integrationService.Sync(ctx, event.Integration.Type, &integrationCommons.SyncOptions{
				InstallationID: event.Integration.InstallationID,
				EntityDetails:  event.EntityDetails,
				Data:           data,
				Credentials:    event.Integration.Credentials,
			}); err != nil {
				log.Printf("failed to push secret with ID %s for %s integration: %s", row.ID, event.Integration.Type, event.Integration.ID)
				log.Println(err)
			}
			wg.Done()
		}(&event)
	}
	wg.Wait()

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully synced secrets",
	})
}

//	Called when a new row is inserted inside the `events` table.
func EventInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row eventCommons.Event
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

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
	organisation, err := organisations.GetByEnvironment(ctx, client, row.EnvID)
	if err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to get organisation to which this event is associated with",
			Error:   err.Error.Error(),
		})
	}

	//	Get the integration to which this event belong to.
	integration, err := integrations.GetService().Get(ctx, client, row.IntegrationID)
	if err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to get integration to which this event is associated with",
			Error:   err.Error.Error(),
		})
	}

	response, err := secrets.GetAll(ctx, client, &secretCommons.GetSecretOptions{
		KeyPath: organisation.ID,
		EnvID:   row.EnvID,
	})
	if err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to get secrets associated with this event",
			Error:   err.Error.Error(),
		})
	}

	data := make(map[string]secretCommons.Payload)

	//	Decrypt the value of every secret.
	for key, payload := range response.Data {

		//	If the secret is of type `ciphertext`,
		//	we will need to decrypt it first.
		if payload.Type == secretCommons.Ciphertext {

			secret, err := secrets.Decrypt(ctx, &secretCommons.DecryptSecretOptions{
				Data: secretCommons.Data{
					Key:     key,
					Payload: payload,
				},
				KeyLocation: organisation.ID,
				EnvID:       row.EnvID,
			})
			if err != nil {
				return c.JSON(http.StatusBadGateway, &APIResponse{
					Code:    http.StatusBadRequest,
					Message: "failed to decrypt value of secret: " + key,
					Error:   err.Error.Error(),
				})
			}

			//	Base64 decode the secret value
			b64Decoded, er := base64.StdEncoding.DecodeString(secret.Data.Plaintext)
			if er != nil {
				return c.JSON(http.StatusBadGateway, &APIResponse{
					Code:    http.StatusBadRequest,
					Message: "failed to base 64 decode the decrypted value of secret: " + key,
					Error:   er.Error(),
				})
			}

			payload.Value = string(b64Decoded)

		} else if payload.Type == secretCommons.Plaintext {

			//	Base64 decode the secret value
			b64Decoded, er := base64.StdEncoding.DecodeString(payload.Value.(string))
			if er != nil {
				return c.JSON(http.StatusBadGateway, &APIResponse{
					Code:    http.StatusBadRequest,
					Message: "failed to base 64 decode the decrypted value of secret: " + key,
					Error:   er.Error(),
				})
			}

			payload.Value = string(b64Decoded)
		}

		data[key] = payload
	}

	//	Get the integration service
	integrationService := integrations.GetService()

	if err := integrationService.Sync(ctx, integration.Type, &integrationCommons.SyncOptions{
		InstallationID: integration.InstallationID,
		EntityDetails:  row.EntityDetails,
		Data:           data,
		Credentials:    integration.Credentials,
	}); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("failed to push secret with ID %s for %s integration: %s", row.ID, integration.Type, row.IntegrationID),
			Error:   err.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully synced secrets",
	})
}

//	Called when a new row is inserted inside the `users` table.
func UserInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var user userCommons.User
	if err := MapToStruct(payload.Event.Data.New, &user); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Create a new `default` organisation for the new user.
	_, err := organisations.CreateWithUserID(ctx, client, &organisations.CreateOptions{
		Name:   fmt.Sprintf("%s's Org", strings.Split(user.Name, "")[0]),
		UserID: user.ID,
	})
	if err != nil {
		return c.JSON(http.StatusExpectationFailed, &APIResponse{
			Code:    http.StatusExpectationFailed,
			Message: "failed to create default organisation",
			Error:   err.Message,
		})
	}

	//	TODO: Shoot a welcome email to the user

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully generated the transit key",
	})
}

//	Called when a new row is inserted inside the `organisations` table.
func OrganisationCreateKey(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row organisations.Organisation
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Generate new transit for this organisation in vault.
	if err := secrets.GenerateKey(ctx, row.ID, secretCommons.GenerateKeyOptions{
		Exportable:           true,
		AllowPlaintextBackup: true,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, &APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to generate transit key",
			Error:   err.Message,
		})
	}

	//	Export the key.
	key, err := secrets.BackupKey(ctx, row.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to export transit key",
			Error:   err.Message,
		})
	}

	//	Get the mailer service.
	service := mail.GetService()

	//	Email the key to owner.
	if err := service.SendKey(ctx, &commons.SendKeyOptions{
		Key:     key.Data.Backup,
		UserID:  row.UserID,
		OrgName: row.Name,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, &APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to email transit key",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully generated the transit key",
	})
}

//	Called when a new row is inserted inside the `organisations` table.
func OrganisationCreateDefaultRoles(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row organisations.Organisation
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Generate default roles for the organisation.
	roles := []permissionCommons.RoleInsertOptions{
		{
			Name: "viewer",
			Permissions: permissionCommons.Permissions{
				Projects: permissionCommons.CRUD{
					Read: true,
				},
			},
		},
		{
			Name: "editor",
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
		},
		{
			Name: "admin",
			Permissions: permissionCommons.Permissions{
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
		},
	}

	//	Initialize the permissions service.
	service := permissions.GetService()

	//	Create the roles.
	for _, item := range roles {

		//	Set the organisation ID in Role
		item.OrgID = row.ID

		if err := service.Insert(permissionCommons.RoleLevelPermission, ctx, client, item); err != nil {
			return c.JSON(http.StatusBadRequest, &APIResponse{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("failed to create role %s for org_id %s", item.Name, item.OrgID),
				Error:   err.Message,
			})
		}
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully generated the transit key and created roles",
	})
}

//	Called when a row is deleted from the `organisations` table.
func OrganisationDeleted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var organisation organisations.Organisation
	if err := MapToStruct(payload.Event.Data.Old, &organisation); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Generate new transit for this organisation in vault.
	if err := secrets.DeleteKey(ctx, organisation.ID); err != nil {
		return c.JSON(http.StatusExpectationFailed, &APIResponse{
			Code:    http.StatusExpectationFailed,
			Message: "failed to delete the transit key",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully deleted the transit key",
	})
}

//	Called when a new row is inserted inside the `projects` table.
func ProjectInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row projects.Project
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Create default environments for this new project.
	envs := []string{"dev", "staging", "qa", "production"}
	for _, item := range envs {
		if _, err := environments.CreateWithUserID(ctx, client, &environments.CreateOptions{
			Name:      item,
			ProjectID: row.ID,
			UserID:    row.UserID,
		}); err != nil {
			return c.JSON(http.StatusBadRequest, &APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to create environment: " + item,
				Error:   err.Error.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully default environments",
	})
}

//	Called when a new row is inserted inside the `invites` table.
func InviteInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row inviteCommons.Invite
	if err := MapToStruct(payload.Event.Data.New, &row); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Get the mailer service.
	service := mail.GetService()

	//	Send the invitation email.
	if err := service.Invite(ctx, &commons.InvitationOptions{
		ID:            row.ID,
		Key:           row.Key,
		ReceiverEmail: row.Email,
		OrgID:         row.OrgID,
		SenderID:      row.UserID,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: err.Message,
			Error:   err.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully sent invitation email to " + row.Email,
	})
}

/* //	Called when a row is inserted/updated/deleted inside the `project_level_permissions` table.
func ProjectLevelPermissions(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var project permissionCommons.ProjectPermissions
	if err := MapToStruct(payload.Event.Data.New, &project); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	incomingPermissions, err := project.GetPermissions()
	if err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal incoming permissions",
			Error:   err.Error(),
		})
	}

	//	Fetch the permissions service
	service := permissions.GetService()

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	switch payload.Event.Op {
	case string(Insert):

		var permissions permissionCommons.Permissions

		//	If the user has been given permission to "write secrets" in the project,
		//	we have to give the user permission to write secrets in every environment of every project.
		if incomingPermissions.SecretsWrite {
			permissions.SecretsWrite = true
		}

		//	If the user has been given permission to "manage permissions" in the project,
		//	we have to give the user permission to manage permissions in every environment of every project.
		if incomingPermissions.PermissionsManage {
			permissions.PermissionsManage = true
		}

		//	Fetch all projects of the organisation
		environments, err := environments.List(ctx, client, &environments.ListOptions{
			ProjectID: project.ProjectID,
		})
		if err != nil {
			return c.JSON(http.StatusBadGateway, &APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to fetch projects for project",
				Error:   err.Message,
			})
		}

		//	Insert permissions for every project
		for _, item := range *environments {
			if err := service.Insert(
				permissionCommons.EnvironmentLevelPermission,
				ctx,
				client,
				permissionCommons.EnvironmentPermissionsInsertOptions{
					EnvID:       item.ID,
					UserID:      project.UserID,
					Permissions: permissions}); err != nil {
				return c.JSON(http.StatusBadGateway, &APIResponse{
					Code:    http.StatusBadRequest,
					Message: "failed to insert permissions for environment: " + item.ID,
					Error:   err.Message,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully inserted environment level permissions",
	})
}

//	Called when a row is inserted/updated/deleted inside the `env_level_permissions` table.
func EnvironmentLevelPermissions(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	return c.JSON(http.StatusBadRequest, &APIResponse{
		Code:    http.StatusBadRequest,
		Message: "un-built event endpoint",
	})
}
*/
