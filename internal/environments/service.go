package environments

import (
	"encoding/json"
	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/events"
	eventCommons "github.com/envsecrets/envsecrets/internal/events/commons"
	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/machinebox/graphql"
)

// Create a new environment
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Environment, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Environment, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

// Get a environment by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Environment, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get a environment by ID
func GetByNameAndProjectID(ctx context.ServiceContext, client *clients.GQLClient, name, project_id string) (*Environment, error) {

	req := graphql.NewRequest(`
	query MyQuery($name: String!, $project_id: uuid!) {
		environments(where: {_and: [{name: {_eq: $name}}, {project_id: {_eq: $project_id}}]}) {
		  id
		}
	  }			
	`)

	req.Var("name", name)
	req.Var("project_id", project_id)

	var response struct {
		Result []Environment `json:"environments"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	if len(response.Result) == 0 {
		return nil, errors.New("environment not found")
	}

	return &response.Result[0], nil
}

// List environments
func List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) (*[]Environment, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update a environment by ID
func Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) (*Environment, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Environment
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Delete a environment by ID
func Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) error {
	return nil
}

// This function syncs the secrets of an environment with it's connected integrations.
// This function assumed that the secrets being supplied are already decrypted.
func Sync(ctx context.ServiceContext, client *clients.GQLClient, options *SyncOptions) error {

	var eventList *eventCommons.Events
	var err error
	if options.IntegrationType != "" {
		eventList, err = events.GetByEnvironmentAndIntegrationType(ctx, client, options.EnvID, options.IntegrationType)
		if err != nil {
			return err
		}
	} else {
		eventList, err = events.GetByEnvironment(ctx, client, options.EnvID)
		if err != nil {
			return err
		}
	}

	if eventList == nil || len(*eventList) == 0 {
		return errors.New("there are no events in this environment to sync this secret with")
	}

	//	Get the integration service
	integrationService := integrations.GetService()
	for _, event := range *eventList {
		if err := integrationService.Sync(ctx, client, &integrations.SyncOptions{
			IntegrationID: event.Integration.ID,
			EventID:       event.ID,
			EntityDetails: event.EntityDetails,
			Data:          options.Secrets,
		}); err != nil {
			return err
		}
	}

	return nil
}
