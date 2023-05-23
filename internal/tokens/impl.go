package tokens

import (
	"encoding/hex"
	"time"

	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/tokens/commons"
	"github.com/envsecrets/envsecrets/internal/tokens/graphql"
)

func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateServiceOptions) (*commons.Token, error) {

	now := time.Now()
	exp := now.Add(options.Expiry)

	//	Generate a random token.
	token, err := globalCommons.GenerateRandomBytes(16)
	if err != nil {
		return nil, err
	}

	//	Hash the value
	hash := globalCommons.SHA256Hash(token)

	result, err := graphql.Create(ctx, client, &commons.CreateGraphQLOptions{
		EnvID:  options.EnvID,
		Name:   options.Name,
		Expiry: exp,
		Hash:   hash,
	})
	if err != nil {
		return nil, err
	}

	result.Hash = hex.EncodeToString(token)
	return result, nil
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Token, error) {
	return graphql.Get(ctx, client, id)
}

func GetByHash(ctx context.ServiceContext, client *clients.GQLClient, hash string) (*commons.Token, error) {
	return graphql.GetByHash(ctx, client, hash)
}
