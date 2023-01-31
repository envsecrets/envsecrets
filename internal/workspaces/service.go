package workspaces

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/machinebox/graphql"
)

//	Create a new workspace
func Create(ctx context.ServiceContext, client *graphql.Client, options *CreateOptions) (*CreateResponse, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!) {
		insert_workspaces(objects: {name: $name}) {
		  returning {
			id
			name
		  }
		}
	  }	  
	`)

	req.Var("name", options.Name)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_workspaces"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []CreateResponse
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	result := resp[0]

	//	Add yourself as the first member of the organization
	if _, err = memberships.Create(ctx, client, &memberships.CreateOptions{
		WorkspaceID: result.ID,
	}); err != nil {
		return nil, err
	}

	return &result, nil
}

//	Get a workspace by ID
func Get(ctx context.ServiceContext, client *graphql.Client, id string) (*Workspace, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		workspaces_by_pk(id: $id) {
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

	returning, err := json.Marshal(response["workspaces_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Workspace
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	List workspaces
func List(ctx context.ServiceContext, client *graphql.Client) (*[]Workspace, error) {

	req := graphql.NewRequest(`
	query MyQuery {
		workspaces {
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

	returning, err := json.Marshal(response["workspaces"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Workspace
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Update a workspace by ID
func Update(ctx context.ServiceContext, client *graphql.Client, id string, options *UpdateOptions) (*Workspace, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $name: String!) {
		update_workspaces_by_pk(pk_columns: {id: $id}, _set: {name: $name}) {
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

	returning, err := json.Marshal(response["update_workspaces_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Workspace
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Delete a workspace by ID
func Delete(ctx context.ServiceContext, client *graphql.Client, id string) error {
	return nil
}
