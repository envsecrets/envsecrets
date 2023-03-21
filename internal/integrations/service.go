package integrations

import (
	"net/url"
	"strings"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/github"
	"github.com/envsecrets/envsecrets/internal/integrations/internal/vercel"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Integration, *errors.Error)
	ListEntities(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, string) (interface{}, *errors.Error)
	Setup(context.ServiceContext, commons.IntegrationType, url.Values) *errors.Error
	Sync(context.ServiceContext, commons.IntegrationType, *commons.SyncOptions) *errors.Error
}

type DefaultIntegrationService struct{}

func (*DefaultIntegrationService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Integration, *errors.Error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultIntegrationService) ListEntities(ctx context.ServiceContext, client *clients.GQLClient, integrationType commons.IntegrationType, integrationID string) (interface{}, *errors.Error) {

	//	Fetch installation ID for integration.
	integration, err := graphql.Get(ctx, client, integrationID)
	if err != nil {
		return nil, err
	}

	switch integrationType {
	case commons.Github:
		return github.ListEntities(ctx, integration)
	case commons.Vercel:
		return nil, nil
	}

	return nil, nil
}

func (*DefaultIntegrationService) Setup(ctx context.ServiceContext, integrationType commons.IntegrationType, params url.Values) *errors.Error {

	//	Extract the Organisation ID and Authorization token from State.
	payload := strings.Split(params.Get("state"), "/")
	if len(payload) != 2 {
		return errors.New(nil, "invalid callback state", errors.ErrorTypeBadRequest, errors.ErrorSourceGithub)
	}
	orgID := payload[0]
	token := payload[1]

	//	Initialize Hasura client with token extract from state parameter
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   string(clients.AuthorizationHeader),
				Value: "Bearer " + token,
			},
		},
	})

	switch integrationType {
	case commons.Github:
		return github.Setup(ctx, client, &github.SetupOptions{
			InstallationID: params.Get("installation_id"),
			SetupAction:    params.Get("setup_action"),
			State:          params.Get("state"),
			OrgID:          orgID,
			Token:          token,
		})
	case commons.Vercel:
		return vercel.Setup(ctx, client, &vercel.SetupOptions{
			ConfigurationID: params.Get("configurationId"),
			Next:            params.Get("next"),
			Source:          params.Get("source"),
			Code:            params.Get("code"),
			State:           params.Get("state"),
			OrgID:           orgID,
			Token:           token,
		})
	}

	return nil
}

func (*DefaultIntegrationService) Sync(ctx context.ServiceContext, integrationType commons.IntegrationType, options *commons.SyncOptions) *errors.Error {

	switch integrationType {
	case commons.Github:
		return github.Sync(ctx, options)
	default:
		return nil
	}
}
