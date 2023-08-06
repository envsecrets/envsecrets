package graphql

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations/commons"
	"github.com/machinebox/graphql"
)

// Create a new organisation
func Create(ctx context.ServiceContext, client *clients.GQLClient, name string) (*commons.Organisation, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!) {
		insert_organisations(objects: {name: $name}) {
		  returning {
			id
			name
		  }
		}
	  }
	`)

	req.Var("name", name)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_organisations"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

// Create a new organisation
func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateOptions) (*commons.Organisation, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $user_id: uuid!) {
		insert_organisations(objects: {name: $name, user_id: $user_id}) {
		  returning {
			id
			name
		  }
		}
	  }
	`)

	req.Var("name", options.Name)
	req.Var("user_id", options.UserID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_organisations"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

// Get a organisation by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
			id
			name
		}
	  }	  
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["organisations_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp commons.Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get an organisation's key copy of the server
func GetServerKeyCopy(ctx context.ServiceContext, client *clients.GQLClient, id string) ([]byte, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
			server_copy
		}
	  }	  
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["organisations_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp commons.Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if resp.ServerKey == "" {
		return nil, errors.New("no key copy found for server")
	}

	//	Base64 decode the server's key copy
	result, err := base64.StdEncoding.DecodeString(resp.ServerKey)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get a organisation by ID
func GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*commons.Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		organisations(where: {projects: {environments: {id: {_eq: $env_id}}}}, limit: 1) {
		  id
		}
	  }	   
	`)

	req.Var("env_id", env_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["organisations"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

// List organisations
func List(ctx context.ServiceContext, client *clients.GQLClient) (*[]commons.Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery {
		organisations {
			id
			name
		}
	  }	  
	`)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["organisations"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update a organisation by ID
func Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *commons.UpdateOptions) (*commons.Organisation, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $name: String!) {
		update_organisations_by_pk(pk_columns: {id: $id}, _set: {name: $name}) {
			id
		  name
		}
	  }	  
	`)

	req.Var("id", id)
	req.Var("name", options.Name)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["update_organisations_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp commons.Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func UpdateInviteLimit(ctx context.ServiceContext, client *clients.GQLClient, options *commons.UpdateInviteLimitOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $limit: Int!) {
		update_organisations(where: {id: {_eq: $id}}, _inc: {invite_limit: $limit}) {
			affected_rows
		  }
	  }			
	`)

	req.Var("id", options.ID)
	req.Var("limit", options.IncrementLimitBy)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["update_organisations"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to update the invite limit")
	}

	return nil
}

func UpdateServerKeyCopy(ctx context.ServiceContext, client *clients.GQLClient, options *commons.UpdateServerKeyCopyOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $key: String!) {
		update_organisations(where: {id: {_eq: $id}}, _set: {server_copy: $key}) {
		  affected_rows
		}
	  }				  
	`)

	req.Var("id", options.OrgID)
	req.Var("key", options.Key)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["update_organisations"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to update the server copy of org's key")
	}

	return nil
}

// Delete a organisation by ID
func Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) error {
	return nil
}
