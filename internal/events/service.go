package events

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations"

	"github.com/machinebox/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*Event, error)
	GetBySecret(context.ServiceContext, *clients.GQLClient, string) (*Events, error)
	GetByEnvironment(context.ServiceContext, *clients.GQLClient, string) (*Events, error)
	GetByEnvironmentAndIntegrationType(context.ServiceContext, *clients.GQLClient, string, integrations.Type) (*Events, error)
	GetByIntegration(context.ServiceContext, *clients.GQLClient, string) (*Events, error)
}

type DefaultService struct{}

func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Event, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		events_by_pk(id: $id) {
		  env_id
		  entity_details
		  integration {
			id
		  }
		}
	  }				
	`)

	req.Var("id", id)

	var response struct {
		Event Event `json:"events_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Event, nil
}

func (*DefaultService) GetBySecret(ctx context.ServiceContext, client *clients.GQLClient, secret_id string) (*Events, error) {

	req := graphql.NewRequest(`
	query MyQuery($secret_id: uuid!) {
		events(where: {environment: {secrets: {id: {_eq: $secret_id}}}}) {
			id
		  env_id
		  entity_details
		  integration {
			id
			installation_id
			type
			credentials
		  }
		}
	  }			  
	`)

	req.Var("secret_id", secret_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["events"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (*DefaultService) GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*Events, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		events(where: {env_id: {_eq: $env_id}}) {
		  id
		  entity_details
		  integration {
			id
			type
		  }
		}
	  }			  
	`)

	req.Var("env_id", env_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["events"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (*DefaultService) GetByEnvironmentAndIntegrationType(ctx context.ServiceContext, client *clients.GQLClient, env_id string, integration_type integrations.Type) (*Events, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!, $integration_type: String!) {
		events(where: {_and: {env_id: {_eq: $env_id}, integration: {type: {_eq: $integration_type}}}}) {
		  id
		  entity_details
		  integration {
			id
		  }
		}
	  }					
	`)

	req.Var("env_id", env_id)
	req.Var("integration_type", integration_type)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["events"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (*DefaultService) GetByIntegration(ctx context.ServiceContext, client *clients.GQLClient, integration_id string) (*Events, error) {

	req := graphql.NewRequest(`
	query MyQuery($integration_id: uuid!) {
		events(where: {integration_id: {_eq: $integration_id}}) {
		  id
		  entity_details
		  integration {
			id
		  }
		}
	  }			  
	`)

	req.Var("integration_id", integration_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["events"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
