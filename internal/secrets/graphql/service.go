package graphql

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/machinebox/graphql"
)

func GetAll(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetOptions) (*commons.GetAllResponse, *errors.Error) {

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {
		return GetAllByVersion(ctx, client, options)
	}

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		secrets(where: {env_id: {_eq: $env_id}}, order_by: {version: desc}, limit: 1) {
		  data
		  version
		}
	  }			
	`)

	req.Var("env_id", options.EnvID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returnedInterface := response["secrets"]
	returned := returnedInterface.([]interface{})

	//	If there are no secrets, return an empty response
	if returnedInterface == nil || len(returned) == 0 {
		return &commons.GetAllResponse{}, nil
	}

	payload := returned[0].(map[string]interface{})
	return &commons.GetAllResponse{
		Data:    payload["data"].(map[string]interface{}),
		Version: int(payload["version"].(float64)),
	}, nil
}

func GetAllByVersion(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetOptions) (*commons.GetAllResponse, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!, $version: Int!) {
		secrets(where: {env_id: {_eq: $env_id}, version: {_eq: $version}}) {
		  data
		  version
		}
	  }			  
	`)

	req.Var("env_id", options.EnvID)
	req.Var("version", options.Version)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returnedInterface := response["secrets"]
	returned := returnedInterface.([]interface{})

	//	If there are no secrets, return an empty response
	if returnedInterface == nil || len(returned) == 0 || returned[0] == nil {
		return &commons.GetAllResponse{}, nil
	}

	payload := returned[0].(map[string]interface{})
	return &commons.GetAllResponse{
		Data:    payload["data"].(map[string]interface{}),
		Version: int(payload["version"].(float64)),
	}, nil
}

func GetByKey(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetOptions) (*commons.Secret, *errors.Error) {

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {
		return GetByKeyByVersion(ctx, client, options)
	}

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!, $key: String!) {
		secrets(order_by: {version: desc}, limit: 1, where: {env_id: {_eq: $env_id}}) {
		  data(path: $key)
		  version
		}
	  }			  
	`)

	req.Var("env_id", options.EnvID)
	req.Var("key", options.Secret.Key)

	var response struct {
		Secrets []map[string]interface{} `json:"secrets"`
	}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	var secret commons.Secret
	secret.Value = response.Secrets[0]["data"]

	return &secret, nil
}

func GetByKeyByVersion(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetOptions) (*commons.Secret, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!, $key: String!, $version: Int!) {
		secrets(where: {env_id: {_eq: $env_id}, version: {_eq: $version}}) {
		  data(path: $key)
		  version
		}
	  }					
	`)

	req.Var("env_id", options.EnvID)
	req.Var("key", options.Secret.Key)
	req.Var("version", options.Version)

	var response struct {
		Secrets []map[string]interface{} `json:"secrets"`
	}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	var secret commons.Secret
	secret.Value = response.Secrets[0]["data"]

	return &secret, nil
}

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetRequestOptions) *errors.Error {

	req := graphql.NewRequest(`
	mutation MyMutation($env_id: uuid!, $data: jsonb!, $version: Int!) {
		insert_secrets(objects: {env_id: $env_id, data: $data, version: $version}) {
		  affected_rows
		}
	  }					
	`)

	//	Fetch the secret of latest version.
	latestEntry, err := GetAll(ctx, client, &commons.GetOptions{
		EnvID: options.EnvID,
	})
	if err != nil {
		return err
	}

	//	We need to create an incremented version.
	incrementBy := 1
	version := latestEntry.Version + incrementBy

	payload := make(map[string]interface{})
	if latestEntry.Data != nil {
		payload = latestEntry.Data
	}

	//	Update our key in the data
	payload[options.Secret.Key] = options.Secret.Value

	//	Set the variables for our GQL query.
	req.Var("env_id", options.EnvID)
	req.Var("version", version)
	req.Var("data", payload)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_secrets"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, "failed to save ciphered secret value", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}
