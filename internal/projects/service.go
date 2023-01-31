package projects

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

//	Create a new workspace
func Create(ctx context.ServiceContext, client *graphql.Client, options *CreateOptions) (*CreateResponse, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $org_id: uuid!) {
		insert_projects(objects: {name: $name, org_id: $org_id}) {
		  returning {
			id
			name
		  }
		}
	  }	  
	`)

	req.Var("name", options.Name)
	req.Var("org_id", options.OrgID)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_projects"].(map[string]interface{})["returning"].([]interface{}))
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

//	Get a workspace by ID
func Get(ctx context.ServiceContext, client *graphql.Client, id string) (*Project, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		projects_by_pk(id: $id) {
			id
			name
			org_id
		}
	  }	  
	`)

	req.Var("id", id)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["projects_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Project
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	List projects
func List(ctx context.ServiceContext, client *graphql.Client, options *ListOptions) (*[]Project, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		projects(where: {org_id: {_eq: $id}}) {
		  id
		  name
		}
	  }	  
	`)

	req.Var("id", options.OrgID)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["projects"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Project
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Update a workspace by ID
func Update(ctx context.ServiceContext, client *graphql.Client, id string, options *UpdateOptions) (*Project, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $name: String!) {
		update_projects_by_pk(pk_columns: {id: $id}, _set: {name: $name}) {
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

	returning, err := json.Marshal(response["update_projects_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Project
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Delete a workspace by ID
func Delete(ctx context.ServiceContext, client *graphql.Client, id string) error {
	return nil
}
