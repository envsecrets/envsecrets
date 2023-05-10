package graphql

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/machinebox/graphql"
)

func Insert(ctx context.ServiceContext, client *clients.GQLClient, options *commons.AddIntegrationOptions) (*commons.Integration, *errors.Error) {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!, $installation_id: String, $type: String!, $credentials: String) {
		insert_integrations(objects: {org_id: $org_id, installation_id: $installation_id, type: $type, credentials: $credentials}) {
		  returning {
			id
		  }
		}
	  }						
	`)

	req.Var("org_id", options.OrgID)
	req.Var("type", options.Type)

	if options.InstallationID != "" {
		req.Var("installation_id", options.InstallationID)
	}
	if options.Credentials != "" {
		req.Var("credentials", options.Credentials)
	}

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returned := response["insert_integrations"].(map[string]interface{})

	returning, err := json.Marshal(returned["returning"].([]interface{}))
	if err != nil {
		return nil, errors.New(err, "failed to marhshal integration into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Integration
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarhshal integration into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
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

func UpdateDetails(ctx context.ServiceContext, client *clients.GQLClient, options *commons.UpdateDetailsOptions) *errors.Error {

	errorMessage := "Failed to update entity details"

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $details: jsonb!) {
		update_events(where: {id: {_eq: $id}}, _set: {entity_details: $details}) {
		  affected_rows
		}
	  }							  
	`)

	req.Var("id", options.ID)
	req.Var("details", options.EntityDetails)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["update_events"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, errorMessage, errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}
