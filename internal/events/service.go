package events

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/events/commons"
	"github.com/machinebox/graphql"
)

func GetBySecret(ctx context.ServiceContext, client *clients.GQLClient, secret_id string) (*commons.Events, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($secret_id: uuid!) {
		events(where: {environment: {secrets: {id: {_eq: $secret_id}}}}) {
		  env_id
		  entity_details
		  integration {
			id
			installation_id
			type
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
		return nil, errors.New(err, "failed to marhshal secrets into json", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp commons.Events
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarhshal secrets into json", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}
