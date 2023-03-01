package triggers

import (
	"encoding/base64"
	"log"
	"net/http"
	"sync"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/events"
	eventCommons "github.com/envsecrets/envsecrets/internal/events/commons"
	"github.com/envsecrets/envsecrets/internal/integrations"
	integrationCommons "github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/permissions"
	permissionCommons "github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/secrets"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
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

	data := make(map[string]interface{})

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

				//	Base64 decode the decrypted value
				b64Decoded, er := base64.StdEncoding.DecodeString(secret.Data.Plaintext)
				if er != nil {
					return c.JSON(http.StatusBadGateway, &APIResponse{
						Code:    http.StatusBadRequest,
						Message: "failed to base 64 decode the decrypted value of secret: " + key,
						Error:   err.Error.Error(),
					})
				}

				data[key] = string(b64Decoded)

			} else if payload.Type == secretCommons.Plaintext {

				data[key] = payload.Value
			}

		}
	}

	//	Get the integration service
	integrationService := integrations.GetService()

	var wg sync.WaitGroup
	for _, event := range *events {
		wg.Add(1)
		go func(event *eventCommons.Event) {
			if err := integrationService.PushSecrets(ctx, event.Integration.Type, &integrationCommons.PushSecretOptions{
				InstallationID: event.Integration.InstallationID,
				EntitySlug:     event.EntitySlug,
				Data:           data,
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

//	Called when a new row is inserted inside the `organisations` table.
func OrganisationInserted(c echo.Context) error {

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
	if err := MapToStruct(payload.Event.Data.New, &organisation); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Initialize a new default context
	ctx := context.NewContext(&context.Config{Type: context.APIContext})

	//	Generate new transit for this organisation in vault.
	if err := secrets.GenerateKey(ctx, organisation.ID, commons.GenerateKeyOptions{
		Exportable: true,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, &APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to generate transit key",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully generated the transit key",
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
	if err := MapToStruct(payload.Event.Data.New, &organisation); err != nil {
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
		return c.JSON(http.StatusInternalServerError, &APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete the transit key",
			Error:   err.Message,
		})
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully deleted the transit key",
	})
}

//	Called when a row is inserted/updated/deleted inside the `org_level_permissions` table.
func OrganisationLevelPermissions(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var organisation permissionCommons.OrgnisationPermissions
	if err := MapToStruct(payload.Event.Data.New, &organisation); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	incomingPermissions, err := organisation.GetPermissions()
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

		//	If the user has been given permission to "manage projects" in the organisation,
		//	we have to give the user permission to manage every environment of every project.
		if incomingPermissions.ProjectsManage {
			permissions.EnvironmentsManage = true
		}

		//	If the user has been given permission to "write secrets" in the organisation,
		//	we have to give the user permission to write secrets in every environment of every project.
		if incomingPermissions.SecretsWrite {
			permissions.SecretsWrite = true
		}

		//	If the user has been given permission to "manage permissions" in the organisation,
		//	we have to give the user permission to manage permissions in every environment of every project.
		if incomingPermissions.PermissionsManage {
			permissions.PermissionsManage = true
		}

		//	Fetch all projects of the organisation
		projects, err := projects.List(ctx, client, &projects.ListOptions{
			OrgID: organisation.OrgID,
		})
		if err != nil {
			return c.JSON(http.StatusBadGateway, &APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to fetch projects for organisation",
				Error:   err.Message,
			})
		}

		//	Insert permissions for every project
		for _, item := range *projects {
			if err := service.Insert(
				permissionCommons.ProjectLevelPermission,
				ctx,
				client,
				permissionCommons.ProjectPermissionsInsertOptions{
					ProjectID:   item.ID,
					UserID:      organisation.UserID,
					Permissions: permissions}); err != nil {
				return c.JSON(http.StatusBadGateway, &APIResponse{
					Code:    http.StatusBadRequest,
					Message: "failed to insert permissions for project: " + item.ID,
					Error:   err.Message,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully inserted project level permissions",
	})
}

//	Called when a row is inserted/updated/deleted inside the `project_level_permissions` table.
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
