package invites

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/invites/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Invite, error)
	Update(context.ServiceContext, *clients.GQLClient, string, *commons.UpdateOptions) (*commons.Invite, error)
}

type DefaultInviteService struct{}

func (*DefaultInviteService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Invite, error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultInviteService) Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *commons.UpdateOptions) (*commons.Invite, error) {
	return graphql.Update(ctx, client, id, options)
}
