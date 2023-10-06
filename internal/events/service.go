package events

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/events/commons"
	"github.com/envsecrets/envsecrets/internal/integrations"

	"github.com/machinebox/graphql"
)

func GetBySecret(ctx context.ServiceContext, client *clients.GQLClient, secret_id string) (*commons.Events, error) {

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
	var resp commons.Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*commons.Events, error) {

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
	var resp commons.Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetByEnvironmentAndIntegrationType(ctx context.ServiceContext, client *clients.GQLClient, env_id string, integration_type integrations.Type) (*commons.Events, error) {

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
	var resp commons.Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetByIntegration(ctx context.ServiceContext, client *clients.GQLClient, integration_id string) (*commons.Events, error) {

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
	var resp commons.Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
