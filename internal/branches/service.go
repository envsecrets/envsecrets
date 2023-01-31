package branches

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

//	Create a new branch
func Create(ctx context.ServiceContext, client *graphql.Client, options *CreateOptions) (*CreateResponse, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $environment_id: uuid!) {
		insert_branches_one(objects: {name: $name, environment_id: $environment_id}) {
			id
			name
		}
	  }	  
	`)

	req.Var("name", options.Name)
	req.Var("environment_id", options.EnvironmentID)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_branches_one"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp CreateResponse
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Get a branch by ID
func Get(ctx context.ServiceContext, client *graphql.Client, id string) (*Branch, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		branches_by_pk(id: $id) {
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

	returning, err := json.Marshal(response["branches_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Branch
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	List branches
func List(ctx context.ServiceContext, client *graphql.Client) (*[]Branch, error) {

	req := graphql.NewRequest(`
	query MyQuery {
		branches {
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

	returning, err := json.Marshal(response["branches"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Branch
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Update a branch by ID
func Update(ctx context.ServiceContext, client *graphql.Client, id string, options *UpdateOptions) (*Branch, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $name: String!) {
		update_branches_by_pk(pk_columns: {id: $id}, _set: {name: $name}) {
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

	returning, err := json.Marshal(response["update_branches_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Branch
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Delete a branch by ID
func Delete(ctx context.ServiceContext, client *graphql.Client, id string) error {
	return nil
}
