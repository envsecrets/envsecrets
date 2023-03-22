package graphql

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/machinebox/graphql"
)

func Insert(ctx context.ServiceContext, client *clients.GQLClient, options *commons.AddIntegrationOptions) *errors.Error {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!, $installation_id: String!, $type: String!, $credentials: jsonb) {
		insert_integrations(objects: {org_id: $org_id, installation_id: $installation_id, type: $type, credentials: $credentials}) {
		  affected_rows
		}
	  }						
	`)

	req.Var("org_id", options.OrgID)
	req.Var("installation_id", options.InstallationID)
	req.Var("type", options.Type)

	if options.Credentials != nil {
		req.Var("credentials", options.Credentials)
	}

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_integrations"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, "failed to insert integration", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Integration, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		integrations_by_pk(id: $id) {
		  id
		  installation_id
		  org_id
		  type
		  user_id
		  credentials
		}
	  }					  
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["integrations_by_pk"])
	if err != nil {
		return nil, errors.New(err, "failed to marhshal integration into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp commons.Integration
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarhshal integration into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

func List(ctx context.ServiceContext, client *clients.GQLClient, options *commons.ListIntegrationFilters) (*commons.Integrations, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($org_id: uuid!, $type: String!) {
		integrations(where: {_and: [{org_id: {_eq: $org_id}}, {type: {_eq: $type}}]}) {
		  id
		  installation_id
		  user_id
		}
	  }						
	`)

	req.Var("org_id", options.OrgID)
	req.Var("type", options.Type)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["integrations"])
	if err != nil {
		return nil, errors.New(err, "failed to marhshal integration into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp commons.Integrations
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarhshal integration into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}
