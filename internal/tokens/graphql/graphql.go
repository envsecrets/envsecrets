package graphql

import (
	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/tokens/commons"
	"github.com/machinebox/graphql"
)

//	Create a new organisation
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateGraphQLOptions) (*commons.Token, *errors.Error) {

	errorMessage := "Failed to create token"

	req := graphql.NewRequest(`
	mutation MyMutation($env_id: uuid!, $expiry: timestamp, $id: uuid!) {
		insert_tokens_one(object: {env_id: $env_id, expiry: $expiry, id: $id}) {
		  id
		}
	  }	  
	`)

	req.Var("id", options.ID)
	req.Var("env_id", options.EnvID)
	if !options.Expiry.IsZero() {
		req.Var("expiry", options.Expiry)
	}

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	var resp commons.Token
	if err := globalCommons.MapToStruct(response["insert_tokens_one"].(map[string]interface{}), &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Fetch a token by it's  ID.
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Token, *errors.Error) {

	errorMessage := "Failed to get the token"

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		tokens_by_pk(id: $id) {
		  env_id
		}
	  }				  
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	var resp commons.Token
	if err := globalCommons.MapToStruct(response["tokens_by_pk"], &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Fetch a token by it's environment ID.
func GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*commons.Token, *errors.Error) {

	errorMessage := "Failed to get the token"

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		tokens(where: {env_id: {_eq: $env_id}}) {
		  env_id
		  expiry
		  id
		}
	  }			
	`)

	req.Var("env_id", env_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	var resp []commons.Token
	if err := globalCommons.MapToStruct(response["tokens"].([]interface{}), &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}
