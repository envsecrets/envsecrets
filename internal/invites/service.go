package invites

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/machinebox/graphql"
)

//	Create a new invite
func Create(ctx context.ServiceContext, client *client.GQLClient, options *CreateOptions) (*CreateResponse, *errors.Error) {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!, $scope: String!, $email: citext!) {
		insert_invites_one(object: {org_id: $org_id, receiver_email: $email, scope: $scope}) {
		  id
		}
	  }	  
	`)

	req.Var("org_id", options.OrgID)
	req.Var("email", options.ReceiverEmail)
	req.Var("scope", options.Scope)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_invites_one"])
	if err != nil {
		return nil, errors.New(err, "failed to marshal json returning response", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp CreateResponse
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarshal json returning response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Get a invite by ID
func Get(ctx context.ServiceContext, client *client.GQLClient, id string) (*Invite, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		invites_by_pk(id: $id) {
			id
		}
	  }	  
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["invites_by_pk"])
	if err != nil {
		return nil, errors.New(err, "failed to marshal json returning response", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Invite
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarshal json returning response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	List invites
func List(ctx context.ServiceContext, client *client.GQLClient, options *ListOptions) (*[]Invite, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($accepted: Boolean) {
		invites(where: {accepted: {_eq: $accepted}}) {
		  id
		  created_at
		  organisation {
			id
			name
		  }
		}
	  }		
	`)

	req.Var("accepted", options.Accepted)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["invites"])
	if err != nil {
		return nil, errors.New(err, "failed to marshal json returning response", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Invite
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarshal json returning response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Update a invite by ID
func Update(ctx context.ServiceContext, client *client.GQLClient, id string, options *UpdateOptions) (*Invite, error) {
	return nil, nil
}

//	Delete a invite by ID
func Delete(ctx context.ServiceContext, client *client.GQLClient, id string) error {
	return nil
}
