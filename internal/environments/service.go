package environments

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

//	Create a new environment
func Create(ctx context.ServiceContext, client *graphql.Client, options *CreateOptions) (*CreateResponse, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $project_id: uuid!) {
		insert_environments(objects: {name: $name, project_id: $project_id}) {
		  returning {
			id
			name
		  }
		}
	  }	  
	`)

	req.Var("name", options.Name)
	req.Var("project_id", options.ProjectID)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_environments"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []CreateResponse
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

//	Get a environment by ID
func Get(ctx context.ServiceContext, client *graphql.Client, id string) (*Environment, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		environments_by_pk(id: $id) {
			id
			name
		}
	  }	  
	`)

	req.Var("id", id)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["environments_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	List environments
func List(ctx context.ServiceContext, client *graphql.Client) (*[]Environment, error) {

	req := graphql.NewRequest(`
	query MyQuery {
		environments {
			id
			name
		}
	  }	  
	`)

	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["environments"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Update a environment by ID
func Update(ctx context.ServiceContext, client *graphql.Client, id string, options *UpdateOptions) (*Environment, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $name: String!) {
		update_environments_by_pk(pk_columns: {id: $id}, _set: {name: $name}) {
			id
		  name
		}
	  }	  
	`)

	req.Var("id", id)
	req.Var("name", options.Name)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["update_environments_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Delete a environment by ID
func Delete(ctx context.ServiceContext, client *graphql.Client, id string) error {
	return nil
}
