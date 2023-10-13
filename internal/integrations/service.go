package integrations

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/asm"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/circleci"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/github"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/gitlab"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/gsm"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/hasura"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/netlify"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/nhost"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/railway"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/supabase"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/vercel"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*Integration, error)
	List(context.ServiceContext, *clients.GQLClient, *ListIntegrationFilters) (*Integrations, error)
	ListEntities(context.ServiceContext, *clients.GQLClient, Type, string, map[string]interface{}) (interface{}, error)
	ListSubEntities(context.ServiceContext, *clients.GQLClient, Type, string, url.Values) (interface{}, error)
	Setup(context.ServiceContext, *clients.GQLClient, Type, *SetupOptions) (*Integration, error)
	Sync(context.ServiceContext, *clients.GQLClient, *SyncOptions) error
}

type DefaultService struct{}

func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Integration, error) {

	result, err := graphql.Get(ctx, client, id)
	if err != nil {
		return nil, err
	}

	return &Integration{
		ID:             result.ID,
		OrgID:          result.OrgID,
		InstallationID: result.InstallationID,
		Type:           Type(result.Type),
		Credentials:    result.Credentials,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
		UserID:         result.UserID,
	}, nil
}

func (*DefaultService) List(ctx context.ServiceContext, client *clients.GQLClient, options *ListIntegrationFilters) (*Integrations, error) {

	var integrations Integrations

	result, err := graphql.List(ctx, client, &graphql.ListIntegrationFilters{
		OrgID: options.OrgID,
		Type:  string(options.Type),
	})
	if err != nil {
		return nil, err
	}

	for _, item := range result {
		integrations = append(integrations, Integration{
			ID:             item.ID,
			OrgID:          item.OrgID,
			InstallationID: item.InstallationID,
			Type:           Type(item.Type),
			Credentials:    item.Credentials,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
			UserID:         item.UserID,
		})
	}

	return &integrations, nil
}

func (*DefaultService) ListEntities(ctx context.ServiceContext, client *clients.GQLClient, integrationType Type, integrationID string, options map[string]interface{}) (interface{}, error) {

	//	Fetch installation ID for integration.
	integration, err := graphql.Get(ctx, client, integrationID)
	if err != nil {
		return nil, err
	}

	//	Decrypt the credentials.
	var credentials map[string]interface{}
	if integration.Credentials != "" {
		payload, err := base64.StdEncoding.DecodeString(integration.Credentials)
		if err != nil {
			return nil, err
		}

		decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
			return nil, err
		}
	}

	switch integrationType {
	case Github:
		return github.ListEntities(ctx, &github.ListOptions{
			InstallationID: integration.InstallationID,
		})
	case Gitlab:
		return gitlab.ListEntities(ctx, &gitlab.ListOptions{
			Credentials:   credentials,
			Type:          gitlab.EntityType(options["type"].(string)),
			OrgID:         integration.OrgID,
			IntegrationID: integration.ID,
		})
	case Vercel:
		return vercel.ListEntities(ctx, &vercel.ListOptions{
			Credentials: credentials,
		})
	case Supabase:
		return supabase.ListEntities(ctx, &supabase.ListOptions{
			Credentials: credentials,
		})
	case Netlify:
		return netlify.ListEntities(ctx, &netlify.ListOptions{
			Credentials: credentials,
		})
	case ASM:
		return asm.ListEntities(ctx, &asm.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	case GSM:
		return gsm.ListEntities(ctx, &gsm.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	case CircleCI:
		return circleci.ListEntities(ctx, &circleci.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	case Railway:
		return railway.ListEntities(ctx, &railway.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	case Hasura:
		return hasura.ListEntities(ctx, &hasura.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
		})
	case Nhost:
		return nhost.ListEntities(ctx, &nhost.ListOptions{
			Credentials: credentials,
		})
	default:
		return nil, errors.New("invalid integration type")
	}
}

func (*DefaultService) ListSubEntities(ctx context.ServiceContext, client *clients.GQLClient, integrationType Type, integrationID string, params url.Values) (interface{}, error) {

	//	Fetch installation ID for integration.
	integration, err := graphql.Get(ctx, client, integrationID)
	if err != nil {
		return nil, err
	}

	//	Decrypt the credentials.
	var credentials map[string]interface{}
	if integration.Credentials != "" {
		payload, err := base64.StdEncoding.DecodeString(integration.Credentials)
		if err != nil {
			return nil, err
		}

		decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
			return nil, err
		}
	}

	switch integrationType {
	case CircleCI:
		return circleci.ListSubEntities(ctx, &circleci.ListOptions{
			Credentials: credentials,
			OrgID:       integration.OrgID,
			OrgSlug:     params.Get("org-slug"),
		})
	default:
		return nil, errors.New("invalid integration type")
	}
}

func (*DefaultService) Setup(ctx context.ServiceContext, client *clients.GQLClient, integrationType Type, options *SetupOptions) (*Integration, error) {

	//	Initialize the options to insert a new row of our integration in database.
	var data struct {
		InstallationID string
		Credentials    map[string]interface{}
	}

	//	Prepare the connection credentials to be saved in our database.
	switch integrationType {
	case Github:

		data.InstallationID = fmt.Sprint(options.Options["installation_id"])

	case Netlify:

		data.Credentials = map[string]interface{}{
			"token": fmt.Sprint(options.Options["token"]),
		}

	case Gitlab:

		credentials, err := gitlab.PrepareCredentials(ctx, &gitlab.PrepareCredentialsOptions{
			Code: fmt.Sprint(options.Options["code"]),
		})
		if err != nil {
			return nil, err
		}

		data.Credentials = credentials

	case Vercel:

		credentials, err := vercel.PrepareCredentials(ctx, &vercel.PrepareCredentialsOptions{
			Code: fmt.Sprint(options.Options["code"]),
		})
		if err != nil {
			return nil, err
		}

		data.Credentials = credentials

	case ASM:

		data.Credentials = map[string]interface{}{
			"role_arn": fmt.Sprint(options.Options["role_arn"]),
			"region":   fmt.Sprint(options.Options["region"]),
		}

	case GSM:

		var keys map[string]interface{}
		if err := json.Unmarshal([]byte(options.Options["keys"].(string)), &keys); err != nil {
			return nil, err
		}

		data.Credentials = keys

	case CircleCI:

		data.Credentials = map[string]interface{}{
			"token": fmt.Sprint(options.Options["token"]),
		}

	case Railway:

		data.Credentials = map[string]interface{}{
			"token": fmt.Sprint(options.Options["token"]),
		}

	case Hasura:

		data.Credentials = map[string]interface{}{
			"token": fmt.Sprint(options.Options["token"]),
		}

	case Nhost:

		data.Credentials = map[string]interface{}{
			"token": fmt.Sprint(options.Options["token"]),
		}

	case Supabase:

		data.Credentials = map[string]interface{}{
			"token": fmt.Sprint(options.Options["token"]),
		}

	default:
		return nil, errors.New("unsupported integration type")
	}

	addOptions := graphql.AddIntegrationOptions{
		OrgID:          options.OrgID,
		InstallationID: data.InstallationID,
		Type:           string(integrationType),
	}

	//	Encrypt the credentials
	if data.Credentials != nil {
		credentials, err := commons.EncryptCredentials(ctx, options.OrgID, data.Credentials)
		if err != nil {
			return nil, err
		}

		addOptions.Credentials = base64.StdEncoding.EncodeToString(credentials)
	}

	//	Create a new record in Hasura.
	result, err := graphql.Insert(ctx, client, &addOptions)
	if err != nil {
		return nil, err
	}

	return &Integration{
		ID: result.ID,
	}, nil
}

func (d *DefaultService) Sync(ctx context.ServiceContext, client *clients.GQLClient, options *SyncOptions) error {

	//	Get the integration to which this event belong to.
	integration, err := d.Get(ctx, client, options.IntegrationID)
	if err != nil {
		return err
	}

	//	Decrypt the credentials.
	var credentials map[string]interface{}
	if integration.Credentials != "" {
		payload, err := base64.StdEncoding.DecodeString(integration.Credentials)
		if err != nil {
			return err
		}

		decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
			return err
		}
	}

	switch Type(integration.Type) {
	case Github:
		return github.Sync(ctx, &github.SyncOptions{
			InstallationID: integration.InstallationID,
			EntityDetails:  options.EntityDetails,
			Data:           options.Data,
		})
	case Gitlab:
		return gitlab.Sync(ctx, &gitlab.SyncOptions{
			Credentials:   credentials,
			EntityDetails: options.EntityDetails,
			Data:          options.Data,
			IntegrationID: options.IntegrationID,
			OrgID:         integration.OrgID,
		})
	case Vercel:
		return vercel.Sync(ctx, &vercel.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case CircleCI:
		return circleci.Sync(ctx, &circleci.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case Railway:
		return railway.Sync(ctx, &railway.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case Hasura:
		return hasura.Sync(ctx, &hasura.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case Nhost:
		return nhost.Sync(ctx, &nhost.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case Supabase:
		return supabase.Sync(ctx, &supabase.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case Netlify:
		return netlify.Sync(ctx, &netlify.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case GSM:
		return gsm.Sync(ctx, &gsm.SyncOptions{
			Credentials:   credentials,
			Data:          options.Data,
			EntityDetails: options.EntityDetails,
		})
	case ASM:
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
			return graphql.UpdateDetails(ctx, gqlClient, &graphql.UpdateDetailsOptions{
				ID:            options.EventID,
				EntityDetails: options.EntityDetails,
			})
		}
		return nil

	default:
		return errors.New("invalid integration type")
	}
}
