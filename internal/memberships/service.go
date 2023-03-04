package memberships

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/machinebox/graphql"
)

//	Create a new membership
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*CreateResponse, *errors.Error) {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!) {
		insert_memberships_one(object: {org_id: $org_id}) {
		  id
		}
	  }	  
	`)

	req.Var("org_id", options.OrgID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_memberships_one"])
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

//	Get a membership by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Membership, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		memberships_by_pk(id: $id) {
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

	returning, err := json.Marshal(response["memberships_by_pk"])
	if err != nil {
		return nil, errors.New(err, "failed to marshal json returning response", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Membership
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarshal json returning response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	List memberships
func List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) (*[]Membership, *errors.Error) {

	req := graphql.NewRequest(`
	query MyQuery($org_id: uuid!) {
		memberships(where: {org_id: {_eq: $org_id}}) {
		  id
		  created_at
		  user {
			email
			displayName
		  }
		}
	  }		
	`)

	req.Var("org_id", options.OrgID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["memberships"])
	if err != nil {
		return nil, errors.New(err, "failed to marshal json returning response", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Membership
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to unmarshal json returning response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Update a membership by ID
func Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) (*Membership, error) {
	return nil, nil
}

//	Delete a membership by ID
func Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) error {
	return nil
}
