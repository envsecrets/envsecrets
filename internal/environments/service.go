package environments

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/machinebox/graphql"
)

//	Create a new environment
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Environment, *errors.Error) {

	errorMessage := "Failed to create the environment"

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

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_environments"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Environment, *errors.Error) {

	errorMessage := "Failed to create the environment"

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $project_id: uuid!, $user_id: uuid) {
		insert_environments(objects: {name: $name, project_id: $project_id, user_id: $user_id}) {
		  returning {
			id
			name
		  }
		}
	  }			
	`)

	req.Var("name", options.Name)
	req.Var("project_id", options.ProjectID)
	if options.UserID != "" {
		req.Var("user_id", options.UserID)
	}

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_environments"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

//	Get a environment by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Environment, *errors.Error) {

	errorMessage := "Failed to fetch the environment"

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		environments_by_pk(id: $id) {
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

	returning, err := json.Marshal(response["environments_by_pk"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	List environments
func List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) (*[]Environment, *errors.Error) {

	errorMessage := "Failed to list the environments"

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		environments(where: {project_id: {_eq: $id}}) {
		  id
		  name
		}
	  }	  
	`)

	req.Var("id", options.ProjectID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["environments"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Update a environment by ID
func Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) (*Environment, *errors.Error) {

	errorMessage := "Failed to update the environment"

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

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["update_environments_by_pk"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Delete a environment by ID
func Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) error {
	return nil
}
