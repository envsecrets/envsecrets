package integrations

import (
	"encoding/base64"
	"encoding/json"
	internalErrors "errors"
	"net/url"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/asm"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/circle"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/github"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/vercel"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Integration, *errors.Error)
	ListEntities(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, string) (interface{}, *errors.Error)
	ListSubEntities(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, string, url.Values) (interface{}, *errors.Error)
	Setup(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, *commons.SetupOptions) (*commons.Integration, *errors.Error)
	Sync(context.ServiceContext, commons.IntegrationType, *commons.SyncOptions) *errors.Error
}

type DefaultIntegrationService struct{}

func (*DefaultIntegrationService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Integration, *errors.Error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultIntegrationService) ListEntities(ctx context.ServiceContext, client *clients.GQLClient, integrationType commons.IntegrationType, integrationID string) (interface{}, *errors.Error) {

	errMessage := "Failed to list entities"

	//	Fetch installation ID for integration.
	integration, err := graphql.Get(ctx, client, integrationID)
	if err != nil {
		return nil, err
	}

	//	Decrypt the credentials.
	payload, er := base64.StdEncoding.DecodeString(integration.Credentials)
	if er != nil {
		return nil, errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
	if err != nil {
		return nil, err
	}

	var credentials map[string]interface{}
	if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
		return nil, errors.New(err, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	switch integrationType {
	case commons.Github:
		return github.ListEntities(ctx, integration)
	case commons.Vercel:
		return vercel.ListEntities(ctx, &vercel.ListOptions{
			Credentials: credentials,
		})
	case commons.ASM:
		return asm.ListEntities(ctx, &asm.ListOptions{
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
	payload, er := base64.StdEncoding.DecodeString(integration.Credentials)
	if er != nil {
		return nil, errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	decryptedCredentials, err := commons.DecryptCredentials(ctx, integration.OrgID, payload)
	if err != nil {
		return nil, err
	}

	var credentials map[string]interface{}
	if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
		return nil, errors.New(err, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
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
			InstallationID: options.Options["installation_id"],
			SetupAction:    options.Options["setup_action"],
			State:          options.Options["state"],
			OrgID:          options.OrgID,
		})
	case commons.Vercel:
		return vercel.Setup(ctx, client, &vercel.SetupOptions{
			ConfigurationID: options.Options["configurationId"],
			Next:            options.Options["next"],
			Source:          options.Options["source"],
			Code:            options.Options["code"],
			State:           options.Options["state"],
			OrgID:           options.OrgID,
		})
	case commons.ASM:
		return asm.Setup(ctx, client, &asm.SetupOptions{
			Region:  options.Options["region"],
			RoleARN: options.Options["role_arn"],
			OrgID:   options.OrgID,
		})
	case commons.CircleCI:
		return circle.Setup(ctx, client, &circle.SetupOptions{
			Token: options.Options["token"],
			OrgID: options.OrgID,
		})
	}

	return nil, errors.New(nil, "invalid integration type", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
}

func (*DefaultIntegrationService) Sync(ctx context.ServiceContext, integrationType commons.IntegrationType, options *commons.SyncOptions) *errors.Error {

	errMessage := "Failed to sync secrets"

	//	Decrypt the credentials.
	payload, er := base64.StdEncoding.DecodeString(options.Credentials)
	if er != nil {
		return errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	decryptedCredentials, err := commons.DecryptCredentials(ctx, options.OrgID, payload)
	if err != nil {
		return err
	}

	var credentials map[string]interface{}
	if err := json.Unmarshal(decryptedCredentials, &credentials); err != nil {
		return errors.New(err, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	switch integrationType {
	case commons.Github:
		return github.Sync(ctx, options)
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
	case commons.ASM:
		resp, err := asm.Sync(ctx, &asm.SyncOptions{
			OrgID:         options.OrgID,
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
