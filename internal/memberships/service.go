package memberships

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

//	Create a new membership
func Create(ctx context.ServiceContext, client *graphql.Client, options *CreateOptions) (*CreateResponse, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!) {
		insert_memberships_one(object: {org_id: $org_id}) {
		  id
		}
	  }	  
	`)

	req.Var("org_id", options.OrgID)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_memberships_one"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp CreateResponse
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Get a membership by ID
func Get(ctx context.ServiceContext, client *graphql.Client, id string) (*Membership, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		memberships_by_pk(id: $id) {
			id
			name
		}
	  }	  
	`)

	req.Var("id", id)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["memberships_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Membership
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	List memberships
func List(ctx context.ServiceContext, client *graphql.Client) (*[]Membership, error) {

	req := graphql.NewRequest(`
	query MyQuery {
		memberships {
			id
			name
		}
	  }	  
	`)

	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["memberships"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Membership
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Update a membership by ID
func Update(ctx context.ServiceContext, client *graphql.Client, id string, options *UpdateOptions) (*Membership, error) {
	return nil, nil
}

//	Delete a membership by ID
func Delete(ctx context.ServiceContext, client *graphql.Client, id string) error {
	return nil
}
