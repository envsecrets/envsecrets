package integrations

import (
	"encoding/base64"
	"encoding/json"
	internalErrors "errors"
	"fmt"
	"net/url"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/asm"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/circle"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/github"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/gitlab"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/gsm"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/netlify"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/supabase"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/vercel"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Integration, *errors.Error)
	ListEntities(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, string, map[string]interface{}) (interface{}, *errors.Error)
	ListSubEntities(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, string, url.Values) (interface{}, *errors.Error)
	Setup(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, *commons.SetupOptions) (*commons.Integration, *errors.Error)
	Sync(context.ServiceContext, *clients.GQLClient, *commons.SyncOptions) *errors.Error
}

type DefaultIntegrationService struct{}

func (*DefaultIntegrationService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Integration, *errors.Error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultIntegrationService) ListEntities(ctx context.ServiceContext, client *clients.GQLClient, integrationType commons.IntegrationType, integrationID string, options map[string]interface{}) (interface{}, *errors.Error) {

	errMessage := "Failed to list entities"

	//	Fetch installation ID for integration.
	integration, err := graphql.Get(ctx, client, integrationID)
	if err != nil {
		return nil, err
	}

	//	Decrypt the credentials.
	var credentials map[string]interface{}
	if integration.Credentials != "" {
		payload, er := base64.StdEncoding.DecodeString(integration.Credentials)
		if er != nil {
			return nil, errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
		}

		decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
			return nil, errors.New(err, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
		}
	}

	switch integrationType {
	case commons.Github:
		return github.ListEntities(ctx, integration)
	case commons.Gitlab:
		return gitlab.ListEntities(ctx, &gitlab.ListOptions{
			Credentials: credentials,
			Type:        gitlab.EntityType(options["type"].(string)),
			Integration: integration,
		})
	case commons.Vercel:
		return vercel.ListEntities(ctx, &vercel.ListOptions{
			Credentials: credentials,
		})
	case commons.Supabase:
		return supabase.ListEntities(ctx, &supabase.ListOptions{
			Credentials: credentials,
		})
	case commons.Netlify:
		return netlify.ListEntities(ctx, &netlify.ListOptions{
			Credentials: credentials,
		})
	case commons.ASM:
		return asm.ListEntities(ctx, &asm.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	case commons.GSM:
		return gsm.ListEntities(ctx, &gsm.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	case commons.CircleCI:
		return circle.ListEntities(ctx, &circle.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	default:
		return nil, errors.New(internalErrors.New("invalid integration type"), errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceHTTP)
	}
}
func (*DefaultIntegrationService) ListSubEntities(ctx context.ServiceContext, client *clients.GQLClient, integrationType commons.IntegrationType, integrationID string, params url.Values) (interface{}, *errors.Error) {

	errMessage := "Failed to list sub entities"

	//	Fetch installation ID for integration.
	integration, err := graphql.Get(ctx, client, integrationID)
	if err != nil {
		return nil, err
	}

	//	Decrypt the credentials.
	var credentials map[string]interface{}
	if integration.Credentials != "" {
		payload, er := base64.StdEncoding.DecodeString(integration.Credentials)
		if er != nil {
			return nil, errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
		}

		decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
			return nil, errors.New(err, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
		}
	}

	switch integrationType {
	case commons.CircleCI:
		return circle.ListSubEntities(ctx, &circle.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
			OrgSlug:     params.Get("org-slug"),
		})
	default:
		return nil, errors.New(internalErrors.New("invalid integration type"), errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceHTTP)
	}
}

func (*DefaultIntegrationService) Setup(ctx context.ServiceContext, client *clients.GQLClient, integrationType commons.IntegrationType, options *commons.SetupOptions) (*commons.Integration, *errors.Error) {

	switch integrationType {
	case commons.Github:
		return github.Setup(ctx, client, &github.SetupOptions{
			InstallationID: fmt.Sprint(options.Options["installation_id"]),
			SetupAction:    fmt.Sprint(options.Options["setup_action"]),
			State:          fmt.Sprint(options.Options["state"]),
			OrgID:          options.OrgID,
		})
	case commons.Netlify:
		return netlify.Setup(ctx, client, &netlify.SetupOptions{
			Token: fmt.Sprint(options.Options["token"]),
			OrgID: options.OrgID,
		})
	case commons.Gitlab:
		return gitlab.Setup(ctx, client, &gitlab.SetupOptions{
			Code:  fmt.Sprint(options.Options["code"]),
			OrgID: options.OrgID,
		})
	case commons.Vercel:
		return vercel.Setup(ctx, client, &vercel.SetupOptions{
			ConfigurationID: fmt.Sprint(options.Options["configurationId"]),
			Next:            fmt.Sprint(options.Options["next"]),
			Source:          fmt.Sprint(options.Options["source"]),
			Code:            fmt.Sprint(options.Options["code"]),
			State:           fmt.Sprint(options.Options["state"]),
			OrgID:           options.OrgID,
		})
	case commons.ASM:
		return asm.Setup(ctx, client, &asm.SetupOptions{
			Region:  fmt.Sprint(options.Options["region"]),
			RoleARN: fmt.Sprint(options.Options["role_arn"]),
			OrgID:   options.OrgID,
		})
	case commons.GSM:

		var keys map[string]interface{}
		if err := json.Unmarshal([]byte(options.Options["keys"].(string)), &keys); err != nil {
			return nil, errors.New(err, "Failed to setup the GSM integration", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
		}

		return gsm.Setup(ctx, client, &gsm.SetupOptions{
			Keys:  keys,
			OrgID: options.OrgID,
		})
	case commons.CircleCI:
		return circle.Setup(ctx, client, &circle.SetupOptions{
			Token: fmt.Sprint(options.Options["token"]),
			OrgID: options.OrgID,
		})
	case commons.Supabase:
		return supabase.Setup(ctx, client, &supabase.SetupOptions{
			Token: fmt.Sprint(options.Options["token"]),
			OrgID: options.OrgID,
		})
	}

	return nil, errors.New(nil, "invalid integration type", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
}

func (*DefaultIntegrationService) Sync(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SyncOptions) *errors.Error {

	errMessage := "Failed to sync secrets"

	//	Get the integration to which this event belong to.
	integration, err := graphql.Get(ctx, client, options.IntegrationID)
	if err != nil {
		return err
	}

	//	Decrypt the credentials.
	var credentials map[string]interface{}
	if integration.Credentials != "" {
		payload, er := base64.StdEncoding.DecodeString(integration.Credentials)
		if er != nil {
			return errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
		}

		decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
			return errors.New(err, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
		}
	}

	switch integration.Type {
	case commons.Github:
		return github.Sync(ctx, &github.SyncOptions{
			InstallationID: integration.InstallationID,
			EntityDetails:  options.EntityDetails,
			Data:           options.Data,
		})
	case commons.Gitlab:
		return gitlab.Sync(ctx, &gitlab.SyncOptions{
			Credentials:   credentials,
			EntityDetails: options.EntityDetails,
			Data:          options.Data,
			IntegrationID: options.IntegrationID,
			OrgID:         integration.OrgID,
		})
	case commons.Vercel:
		return vercel.Sync(ctx, &vercel.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case commons.CircleCI:
		return circle.Sync(ctx, &circle.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case commons.Supabase:
		return supabase.Sync(ctx, &supabase.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case commons.Netlify:
		return netlify.Sync(ctx, &netlify.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case commons.GSM:
		return gsm.Sync(ctx, &gsm.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case commons.ASM:
		resp, err := asm.Sync(ctx, &asm.SyncOptions{
			OrgID:         integration.OrgID,
			Data:          options.Data,
			Credentials:   credentials,
			EntityDetails: options.EntityDetails,
		})
		if err != nil {
			return err
		}

		if resp != nil {
			options.EntityDetails["secret_arn"] = resp.ARN

			//	Save the ARN of created secret in event's entity_details.
			gqlClient := clients.NewGQLClient(&clients.GQLConfig{
				Type: clients.HasuraClientType,
				Headers: []clients.Header{
					clients.XHasuraAdminSecretHeader,
				},
			})
			return graphql.UpdateDetails(ctx, gqlClient, &commons.UpdateDetailsOptions{
				ID:            options.EventID,
				EntityDetails: options.EntityDetails,
			})
		}
		return nil

	default:
		return errors.New(internalErrors.New(errMessage), errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceHTTP)
	}
}
