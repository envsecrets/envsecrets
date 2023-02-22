package integrations

import (
	"net/url"
	"strings"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/github"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
)

type Service interface {
	ListEntities(context.ServiceContext, *clients.GQLClient, commons.IntegrationType, string) (*commons.Entities, *errors.Error)
	Setup(context.ServiceContext, commons.IntegrationType, url.Values) *errors.Error
	//	PushSecrets(commons.IntegrationType, interface{}) *errors.Error
}

type DefaultIntegrationService struct{}

func (*DefaultIntegrationService) ListEntities(ctx context.ServiceContext, client *clients.GQLClient, integrationType commons.IntegrationType, integrationID string) (*commons.Entities, *errors.Error) {

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

	switch integrationType {
	case commons.Github:

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

		return github.Setup(ctx, client, &github.SetupOptions{
			InstallationID: params.Get("installation_id"),
			SetupAction:    params.Get("setup_action"),
			State:          params.Get("state"),
			OrgID:          orgID,
			Token:          token,
		})
	case commons.Vercel:
		return nil
	}

	return nil
}
