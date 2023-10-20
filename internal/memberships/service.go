package memberships

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

// Create a new membership
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!, $role_id: uuid!, $key: String!) {
		insert_org_has_user(objects: {key: $key, role_id: $role_id, org_id: $org_id}) {
		  affected_rows
		}
	  }	  
	`)

	req.Var("org_id", options.OrgID)
	req.Var("role_id", options.RoleID)
	req.Var("key", options.Key)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_org_has_user"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to create membership")
	}

	return nil
}

// Create a new membership
func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($user_id: uuid!, $org_id: uuid!, $role_id: uuid!, $key: String!) {
		insert_org_has_user(objects: {key: $key, role_id: $role_id, org_id: $org_id, user_id: $user_id}) {
		  affected_rows
		}
	  }	  
	`)

	req.Var("user_id", options.UserID)
	req.Var("org_id", options.OrgID)
	req.Var("role_id", options.RoleID)
	req.Var("key", options.Key)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_org_has_user"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to create membership")
	}

	return nil
}

// Get a membership by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Membership, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		org_has_user_by_pk(id: $id) {
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

	returning, err := json.Marshal(response["org_has_user_by_pk"])
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

// Get a membership by ID
func GetKey(ctx context.ServiceContext, client *clients.GQLClient, options *GetKeyOptions) ([]byte, error) {

	req := graphql.NewRequest(`
	query MyQuery($user_id: uuid!, $org_id: uuid!) {
		org_has_user(where: {org_id: {_eq: $org_id}, user_id: {_eq: $user_id}}) {
		  key
		}
	  }			
	`)

	req.Var("user_id", options.UserID)
	req.Var("org_id", options.OrgID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["org_has_user"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Membership
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	result, err := base64.StdEncoding.DecodeString(resp[0].Key)
	if err != nil {
		return nil, err
	}

	return result, nil
}
