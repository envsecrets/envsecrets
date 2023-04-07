package invites

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/invites/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Invite, *errors.Error)
	Update(context.ServiceContext, *clients.GQLClient, string, *commons.UpdateOptions) (*commons.Invite, *errors.Error)
}

type DefaultInviteService struct{}

func (*DefaultInviteService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Invite, *errors.Error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultInviteService) Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *commons.UpdateOptions) (*commons.Invite, *errors.Error) {
	return graphql.Update(ctx, client, id, options)
}
