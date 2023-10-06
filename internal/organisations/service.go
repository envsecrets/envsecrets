package organisations

import (
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*Organisation, error)
	GetByEnvironment(context.ServiceContext, *clients.GQLClient, string) (*Organisation, error)
	Create(context.ServiceContext, *clients.GQLClient, *CreateOptions) (*Organisation, error)
	List(context.ServiceContext, *clients.GQLClient) (*[]Organisation, error)
	UpdateInviteLimit(context.ServiceContext, *clients.GQLClient, *UpdateInviteLimitOptions) error
}

type DefaultOrganisationService struct{}

func (*DefaultOrganisationService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
			id
			name
		}
	  }	  
	`)

	req.Var("id", id)

	var response struct {
		Organisation *Organisation `json:"organisations_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return response.Organisation, nil
}

func (*DefaultOrganisationService) GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		organisations(where: {projects: {environments: {id: {_eq: $env_id}}}}, limit: 1) {
		  id
		}
	  }	   
	`)

	req.Var("env_id", env_id)

	var response struct {
		Organisations []Organisation `json:"organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Organisations[0], nil
}

func (*DefaultOrganisationService) List(ctx context.ServiceContext, client *clients.GQLClient) (*[]Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		organisations(where: {projects: {environments: {id: {_eq: $env_id}}}}, limit: 1) {
		  id
		}
	  }	   
	`)

	var response struct {
		Organisations []Organisation `json:"organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Organisations, nil
}

func (*DefaultOrganisationService) Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Organisation, error) {

	if options.UserID != "" {
		return createWithUserID(ctx, client, options)
	}

	return create(ctx, client, options.Name)
}

func (*DefaultOrganisationService) UpdateInviteLimit(ctx context.ServiceContext, client *clients.GQLClient, options *UpdateInviteLimitOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $limit: Int!) {
		update_organisations(where: {id: {_eq: $id}}, _inc: {invite_limit: $limit}) {
			affected_rows
		  }
	  }			
	`)

	req.Var("id", options.ID)
	req.Var("limit", options.IncrementLimitBy)

	var response struct {
		Result struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"update_organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	if response.Result.AffectedRows == 0 {
		return fmt.Errorf("failed to update the invite limit")
	}

	return nil
}

//
//	--- GraphQL ---
//

// Create a new organisation
func create(ctx context.ServiceContext, client *clients.GQLClient, name string) (*Organisation, error) {

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

	var response struct {
		Result struct {
			Organisations []Organisation `json:"returning"`
		} `json:"insert_organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Result.Organisations[0], nil
}

// Create a new organisation
func createWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Organisation, error) {

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

	var response struct {
		Result struct {
			Organisations []Organisation `json:"returning"`
		} `json:"insert_organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Result.Organisations[0], nil
}
