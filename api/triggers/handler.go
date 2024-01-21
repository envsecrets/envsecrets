package triggers

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/invites"
	"github.com/envsecrets/envsecrets/internal/mail"
	"github.com/envsecrets/envsecrets/internal/mail/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/subscriptions"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/labstack/echo/v4"
)

const (
	KEY_BYTES = 32
)

// Called when a new row is inserted inside the `secrets` table.
func SecretDeleteLegacy(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload clients.HasuraTriggerPayload
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
	orgSubscriptions, err := subscriptions.GetService().List(ctx, client, &subscriptions.ListOptions{OrgID: organisation.ID})
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
	var payload clients.HasuraTriggerPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var user users.User
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

// Called when a new row is inserted inside the `projects` table.
func ProjectInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload clients.HasuraTriggerPayload
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
	var payload clients.HasuraTriggerPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, &clients.APIResponse{
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var row invites.Invite
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
