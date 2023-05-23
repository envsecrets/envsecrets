package events

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/events/commons"
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
