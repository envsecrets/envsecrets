package tokens

import (
	"time"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/tokens/commons"
	"github.com/envsecrets/envsecrets/internal/tokens/graphql"
)

func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateServiceOptions) (*commons.Token, *errors.Error) {

	now := time.Now()
	exp := now.Add(options.Expiry)

	return graphql.Create(ctx, client, &commons.CreateGraphQLOptions{
		EnvID:  options.EnvID,
		Name:   options.Name,
		Expiry: exp,
	})
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Token, *errors.Error) {
	return graphql.Get(ctx, client, id)
}

func GetByHash(ctx context.ServiceContext, client *clients.GQLClient, hash string) (*commons.Token, *errors.Error) {
	return graphql.GetByHash(ctx, client, hash)
}
