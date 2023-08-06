package organisations

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations/commons"
	"github.com/envsecrets/envsecrets/internal/organisations/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Organisation, error)
	GetByEnvironment(context.ServiceContext, *clients.GQLClient, string) (*commons.Organisation, error)
	Create(context.ServiceContext, *clients.GQLClient, *commons.CreateOptions) (*commons.Organisation, error)
	List(context.ServiceContext, *clients.GQLClient) (*[]commons.Organisation, error)
	UpdateInviteLimit(context.ServiceContext, *clients.GQLClient, *commons.UpdateInviteLimitOptions) error
	GetServerKeyCopy(context.ServiceContext, *clients.GQLClient, string) ([]byte, error)
}

type DefaultOrganisationService struct{}

func (*DefaultOrganisationService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Organisation, error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultOrganisationService) GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*commons.Organisation, error) {
	return graphql.GetByEnvironment(ctx, client, env_id)
}

func (*DefaultOrganisationService) GetServerKeyCopy(ctx context.ServiceContext, client *clients.GQLClient, id string) ([]byte, error) {
	return graphql.GetServerKeyCopy(ctx, client, id)
}

func (*DefaultOrganisationService) List(ctx context.ServiceContext, client *clients.GQLClient) (*[]commons.Organisation, error) {
	return graphql.List(ctx, client)
}

func (*DefaultOrganisationService) Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateOptions) (*commons.Organisation, error) {

	if options.UserID != "" {
		return graphql.CreateWithUserID(ctx, client, options)
	}

	return graphql.Create(ctx, client, options.Name)
}

func (*DefaultOrganisationService) UpdateInviteLimit(ctx context.ServiceContext, client *clients.GQLClient, options *commons.UpdateInviteLimitOptions) error {
	return graphql.UpdateInviteLimit(ctx, client, options)
}
