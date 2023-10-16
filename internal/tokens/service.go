package tokens

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/utils"
	"github.com/machinebox/graphql"
)

type Service interface {
	Create(context.ServiceContext, *clients.GQLClient, *CreateOptions) ([]byte, error)
	Get(context.ServiceContext, *clients.GQLClient, string) (*Token, error)
	GetByHash(context.ServiceContext, *clients.GQLClient, string) (*Token, error)
	List(context.ServiceContext, *clients.GQLClient, *ListOptions) ([]*Token, error)
	Decrypt(context.ServiceContext, *clients.GQLClient, []byte, []byte) ([]byte, error)
}

type DefaultService struct{}

func (*DefaultService) Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) ([]byte, error) {

	now := time.Now()
	exp := now.Add(options.Expiry)

	//	Generate a symmetric key for cryptographic operations in this organisation.
	keyBytes, err := utils.GenerateRandomBytes(KEY_BYTES)
	if err != nil {
		return nil, err
	}

	//	Encrypt the org key using newly generated symmetric key
	var key [32]byte
	copy(key[:], keyBytes)
	token, err := keys.SealSymmetrically(options.OrgKey, key)
	if err != nil {
		return nil, err
	}

	//	Hash the token to store it in our DB.
	hash := utils.SHA256Hash(token)

	if _, err := create(ctx, client, &CreateGraphQLOptions{
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

func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Token, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		tokens_by_pk(id: $id) {
		  key
		}
	  }				  
	`)

	req.Var("id", id)

	var response struct {
		Token Token `json:"tokens_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Token, nil
}

func (*DefaultService) GetByHash(ctx context.ServiceContext, client *clients.GQLClient, hash string) (*Token, error) {

	req := graphql.NewRequest(`
	query MyQuery($hash: String!) {
		tokens(where: {hash: {_eq: $hash}}) {
		  env_id
		  expiry
		  key
		}
	  }			
	`)

	req.Var("hash", hash)

	var response struct {
		Tokens []Token `json:"tokens"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	if len(response.Tokens) == 0 {
		return nil, fmt.Errorf("no tokens found")
	}

	return &response.Tokens[0], nil
}

func (*DefaultService) List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) ([]*Token, error) {

	req := graphql.NewRequest(`
	query MyQuery($where: tokens_bool_exp) {
		tokens(where: $where) {
		  id
		  name
		}
	  }	  
	`)

	req.Var("where", options)

	var response struct {
		Tokens []*Token `json:"tokens"`
	}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return response.Tokens, nil
}

func (*DefaultService) Decrypt(ctx context.ServiceContext, client *clients.GQLClient, token, keyBytes []byte) ([]byte, error) {

	var key [32]byte
	copy(key[:], keyBytes)
	decrypted, err := keys.OpenSymmetrically(token, key)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

//
//	--- GraphQL ---
//

func create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateGraphQLOptions) (*Token, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $key: String!, $hash: String!, $env_id: uuid!, $expiry: timestamptz) {
		insert_tokens_one(object: {name: $name, key: $key, hash: $hash, env_id: $env_id, expiry: $expiry}) {
		  id
		}
	  }
	`)

	req.Var("env_id", options.EnvID)
	req.Var("name", options.Name)
	req.Var("key", base64.StdEncoding.EncodeToString(options.Key))
	req.Var("hash", options.Hash)
	if !options.Expiry.IsZero() {
		req.Var("expiry", options.Expiry)
	}

	var response struct {
		Token Token `json:"insert_tokens_one"`
	}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Token, nil
}
