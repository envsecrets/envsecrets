package tokens

import (
	"encoding/base64"
	"errors"
	"time"

	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/tokens/commons"
	"github.com/envsecrets/envsecrets/internal/tokens/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Token, error)
	Create(context.ServiceContext, *clients.GQLClient, *commons.CreateOptions) ([]byte, error)
	Decrypt(context.ServiceContext, *clients.GQLClient, []byte) (*commons.DecryptResponse, error)
}

type DefaultTokenService struct{}

func (*DefaultTokenService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Token, error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultTokenService) Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateOptions) ([]byte, error) {

	now := time.Now()
	exp := now.Add(options.Expiry)

	//	Generate a symmetric key for cryptographic operations in this organisation.
	keyBytes, err := globalCommons.GenerateRandomBytes(commons.KEY_BYTES)
	if err != nil {
		return nil, err
	}

	//	Encrypt the org key using newly generated symmetric key
	var key [32]byte
	copy(key[:], keyBytes)
	token, err := keys.SealAsymmetricallyAnonymous(options.OrgKey, key)
	if err != nil {
		return nil, err
	}

	//	Hash the token to store it in our DB.
	hash := globalCommons.SHA256Hash(token)

	if _, err := graphql.Create(ctx, client, &commons.CreateGraphQLOptions{
		EnvID:  options.EnvID,
		Name:   options.Name,
		Expiry: exp,
		Key:    keyBytes,
		Hash:   hash,
	}); err != nil {
		return nil, err
	}

	return token, nil
}

func (*DefaultTokenService) Decrypt(ctx context.ServiceContext, client *clients.GQLClient, token []byte) (*commons.DecryptResponse, error) {

	//	Hash the token to fetch it from database.
	hash := globalCommons.SHA256Hash(token)
	row, err := graphql.GetByHash(ctx, client, hash)
	if err != nil {
		return nil, err
	}

	//	Parse the token expiry
	now := time.Now()
	expired := now.After(row.Expiry)
	if expired {
		return nil, errors.New("token expired")
	}

	//	Base64 decode the key.
	keyBytes, err := base64.StdEncoding.DecodeString(row.Key)
	if err != nil {
		return nil, err
	}

	//	Encrypt the org key using newly generated symmetric key
	var key [32]byte
	copy(key[:], keyBytes)
	decrypted, err := keys.OpenSymmetrically(token, key)
	if err != nil {
		return nil, err
	}

	return &commons.DecryptResponse{
		OrgKey: decrypted,
		EnvID:  row.EnvID,
		Expiry: row.Expiry,
		Name:   row.Name,
	}, nil
}
