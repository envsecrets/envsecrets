package graphql

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/machinebox/graphql"
)

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *GetOptions) (*commons.Secret, error) {

	req := options.NewRequest()

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["secrets"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Secret
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf(string(clients.ErrorTypeRecordNotFound))
	}

	item := resp[0]
	if item.Data == nil {
		return nil, fmt.Errorf(string(clients.ErrorTypeRecordNotFound))
	}
	item.MarkEncoded()
	if options.Key != "" {
		item.ChangeKey(commons.TEMP_KEY_NAME, options.Key)
	}
	return &item, nil
}

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *SetOptions) (*commons.Secret, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($env_id: uuid!, $data: jsonb!, $version: Int) {
		insert_secrets(objects: {env_id: $env_id, data: $data, version: $version}) {
		  returning {
			version
		  }
		}
	  }					
	`)

	//	Set the variables for our GQL query.
	req.Var("env_id", options.EnvID)
	req.Var("data", options.Data)
	if options.Version != nil {
		req.Var("version", options.Version)
	}

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returned := response["insert_secrets"].(map[string]interface{})["returning"].([]interface{})
	if len(returned) == 0 {
		return nil, fmt.Errorf(string(clients.ErrorTypeRecordNotFound))
	}

	returning, err := json.Marshal(returned)
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Secret
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf(string(clients.ErrorTypeRecordNotFound))
	}

	item := resp[0]
	item.MarkEncoded()
	return &item, nil
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *DeleteOptions) error {

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
