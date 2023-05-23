package graphql

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/machinebox/graphql"
)

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.GetResponse, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		secrets(where: {env_id: {_eq: $env_id}}, order_by: {version: desc}) {
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

	returning, err := json.Marshal(response["secrets"])
	if err != nil {
		return nil, errors.New(err, "failed to marhshal secrets into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Row
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarhshal secrets into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	var payload commons.GetResponse

	if len(resp) > 0 {
		payload.Secrets = resp[0].Data
		payload.Version = &resp[0].Version
	}

	return &payload, nil
}

func GetByVersion(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.GetResponse, *errors.Error) {

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

	returning, err := json.Marshal(response["secrets"])
	if err != nil {
		return nil, errors.New(err, "failed to marhshal secrets into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Row
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarhshal secrets into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	var payload commons.GetResponse

	if len(resp) > 0 {
		payload.Secrets = resp[0].Data
		payload.Version = &resp[0].Version
	}

	return &payload, nil
}

func GetByKey(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Row, *errors.Error) {

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
		return nil, errors.New(err, "failed to marshal secrets into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []struct {
		Version int             `json:"version,omitempty" graphql:"version,omitempty"`
		Data    commons.Payload `json:"data,omitempty" graphql:"data,omitempty"`
	}
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarshal secrets into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	if len(resp) == 0 {
		return nil, errors.New(nil, "no record found", errors.ErrorTypeBadRequest, errors.ErrorSourceGraphQL)
	}

	index := 0
	secrets := make(commons.Secrets)
	secrets.Set(options.Key, resp[index].Data)
	return &commons.Row{
		Version: resp[index].Version,
		Data:    secrets,
	}, nil
}

func GetByKeyByVersion(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Row, *errors.Error) {

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
		return nil, errors.New(err, "failed to marshal secrets into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []struct {
		Version int             `json:"version,omitempty" graphql:"version,omitempty"`
		Data    commons.Payload `json:"data,omitempty" graphql:"data,omitempty"`
	}
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarshal secrets into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	if len(resp) == 0 {
		return nil, errors.New(nil, "no record found", errors.ErrorTypeBadRequest, errors.ErrorSourceGraphQL)
	}

	index := 0
	secrets := make(commons.Secrets)
	secrets.Set(options.Key, resp[index].Data)
	return &commons.Row{
		Version: resp[index].Version,
		Data:    secrets,
	}, nil
}

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetSecretOptions) (*commons.Row, *errors.Error) {

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

	if latestEntry.Secrets == nil {
		latestEntry.Secrets = make(commons.Secrets)
	}

	latestEntry.Secrets.Overwrite(options.Secrets)

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
	req.Var("data", latestEntry.Secrets)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returned := response["insert_secrets"].(map[string]interface{})["returning"].([]interface{})

	if len(returned) == 0 {
		return nil, errors.New(nil, "failed to save ciphered secret value", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	data, er := json.Marshal(returned[0])
	if er != nil {
		return nil, errors.New(er, "failed to marshal returned secrets array", errors.ErrorTypeBadResponse, errors.ErrorSourceGraphQL)
	}

	var secret commons.Row
	if err := json.Unmarshal(data, &secret); err != nil {
		return nil, errors.New(err, "failed to unmarshal created secrets", errors.ErrorTypeBadResponse, errors.ErrorSourceGraphQL)
	}

	return &secret, nil
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) *errors.Error {

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

	data := make(commons.Secrets)
	for key, payload := range latestEntry.Secrets {
		data.Set(key, payload)
	}

	//	Delete our key=value pair.
	delete(data, options.Key)

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
		return errors.New(nil, "failed to save ciphered secret value", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}

func Cleanup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CleanupSecretOptions) *errors.Error {

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
		return errors.New(nil, "Failed to cleanup secrets", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}
