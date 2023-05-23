package graphql

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/internal/payload"
	"github.com/machinebox/graphql"
)

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Secret, error) {

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

	returning, err := json.Marshal(response["secrets"].([]interface{})[0])
	if err != nil {
		return nil, err
	}

	return commons.ParseAndInitialize(returning)
}

func GetByVersion(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Secret, error) {

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

	returning, err := json.Marshal(response["secrets"].([]interface{})[0])
	if err != nil {
		return nil, err
	}

	result, err := commons.ParseAndInitialize(returning)
	if err != nil {
		return nil, err
	}

	if result.IsEmpty() {
		return nil, fmt.Errorf("no record found")
	}

	return result, nil
}

func GetByKey(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Secret, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!, $key: String!) {
		secrets(order_by: {version: desc}, limit: 1, where: {env_id: {_eq: $env_id}}) {
		  data(path: $key)
		  version
		}
	  }			  
	`)

	req.Var("env_id", options.EnvID)
	req.Var("key", options.Key)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["secrets"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []struct {
		Version *int            `json:"version,omitempty"`
		Data    payload.Payload `json:"data,omitempty"`
	}
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("no record found")
	}

	item := resp[0]

	if item.Data.IsEmpty() {
		return nil, errors.New("no record found")
	}

	version := item.Version
	secret := commons.New()
	secret.Version = version
	secret.Set(options.Key, &item.Data)
	secret.MarkEncoded()
	return secret, nil
}

func GetByKeyByVersion(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Secret, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!, $key: String!, $version: Int!) {
		secrets(limit: 1, where: {env_id: {_eq: $env_id}, version: {_eq: $version}}) {
		  data(path: $key)
		  version
		}
	  }
	`)

	req.Var("env_id", options.EnvID)
	req.Var("key", options.Key)
	req.Var("version", options.Version)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["secrets"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []struct {
		Version *int            `json:"version,omitempty"`
		Data    payload.Payload `json:"data,omitempty"`
	}
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("no record found")
	}

	item := resp[0]

	if item.Data.IsEmpty() {
		return nil, errors.New("no record found")
	}

	version := item.Version
	secret := commons.New()
	secret.Version = version
	secret.Set(options.Key, &item.Data)
	secret.MarkEncoded()
	return secret, nil
}

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetSecretOptions) (*commons.Secret, error) {

	//	Fetch the secret of latest version.
	latestEntry, err := Get(ctx, client, &commons.GetSecretOptions{
		EnvID: options.EnvID,
	})
	if err != nil {
		return nil, err
	}

	//	We need to create an incremented version.
	incrementBy := 1
	version := incrementBy
	if latestEntry.Version != nil {
		version += *latestEntry.Version
	}

	latestEntry.Overwrite(options.Data)

	req := graphql.NewRequest(`
	mutation MyMutation($env_id: uuid!, $data: jsonb!, $version: Int!) {
		insert_secrets(objects: {env_id: $env_id, data: $data, version: $version}) {
		  returning {
			version
		  }
		}
	  }					
	`)

	//	Set the variables for our GQL query.
	req.Var("env_id", options.EnvID)
	req.Var("version", version)
	req.Var("data", latestEntry.GetMap())

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returned := response["insert_secrets"].(map[string]interface{})["returning"].([]interface{})

	if len(returned) == 0 {
		return nil, errors.New("no rows affected")
	}

	data, err := json.Marshal(returned[0])
	if err != nil {
		return nil, err
	}

	return commons.ParseAndInitialize(data)
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) error {

	//	Fetch the secret of latest version.
	latestEntry, err := Get(ctx, client, &commons.GetSecretOptions{
		EnvID: options.EnvID,
	})
	if err != nil {
		return err
	}

	//	We need to create an incremented version.
	incrementBy := 1
	version := incrementBy
	if latestEntry.Version != nil {
		version += *latestEntry.Version
	}

	data := commons.New()
	for key, payload := range latestEntry.Data {
		data.Set(key, payload)
	}

	//	Delete our key=value pair.
	data.Delete(options.Key)

	req := graphql.NewRequest(`
	mutation MyMutation($env_id: uuid!, $data: jsonb!, $version: Int!) {
		insert_secrets(objects: {env_id: $env_id, data: $data, version: $version}) {
		  affected_rows
		}
	  }					
	`)

	//	Set the variables for our GQL query.
	req.Var("env_id", options.EnvID)
	req.Var("version", version)
	req.Var("data", data)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_secrets"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func Cleanup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CleanupSecretOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($version: Int!, $env_id: uuid!) {
		delete_secrets(where: {version: {_lte: $version}, env_id: {_eq: $env_id}}) {
		  affected_rows
		}
	  }					  
	`)

	//	Set the variables for our GQL query.
	req.Var("env_id", options.EnvID)
	req.Var("version", options.Version)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["delete_secrets"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
