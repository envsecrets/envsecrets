package projects

import (
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*Project, error)
	Create(context.ServiceContext, *clients.GQLClient, *CreateOptions) (*Project, error)
	List(context.ServiceContext, *clients.GQLClient, *ListOptions) ([]*Project, error)
	Update(context.ServiceContext, *clients.GQLClient, string, *UpdateOptions) (*Project, error)
	Delete(context.ServiceContext, *clients.GQLClient, string) error
}

type DefaultService struct{}

// Get a project by ID
func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Project, error) {

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

	var response struct {
		Project Project `json:"projects_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Project, nil
}

// Create a new project
func (*DefaultService) Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Project, error) {

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

	var response struct {
		Query struct {
			Returning []Project `json:"returning"`
		} `json:"insert_projects"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	if len(response.Query.Returning) == 0 {
		return nil, fmt.Errorf("no project created")
	}

	return &response.Query.Returning[0], nil
}

// List projects
func (*DefaultService) List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) ([]*Project, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		projects(where: {org_id: {_eq: $id}}) {
		  id
		  name
		}
	  }	  
	`)

	req.Var("id", options.OrgID)

	var response struct {
		Projects []*Project `json:"projects"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return response.Projects, nil
}

// Update a project by ID
func (*DefaultService) Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) (*Project, error) {

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

	var response struct {
		UpdateProject Project `json:"update_projects_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.UpdateProject, nil
}

// Delete a project by ID
func (*DefaultService) Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) error {
	return nil
}
