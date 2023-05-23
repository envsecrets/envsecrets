package roles

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

// Insert new permissions.
func Insert(ctx context.ServiceContext, client *clients.GQLClient, options *RoleInsertOptions) (*Role, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $org_id: uuid!, $permissions: jsonb!) {
		insert_roles(objects: {name: $name, org_id: $org_id, permissions: $permissions}) {
		  returning {
			id
		  }
		}
	  }	   
	`)

	req.Var("org_id", options.OrgID)
	req.Var("name", options.Name)
	req.Var("permissions", options.Permissions)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_roles"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Role
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

// Fetches all the roles in an organisation.
func GetRolesByOrgID(ctx context.ServiceContext, client *clients.GQLClient, org_id string) (*[]Role, error) {

	req := graphql.NewRequest(`
	query MyQuery($org_id: uuid!) {
		roles(where: {org_id: {_eq: $org_id}}) {
		  name
		  id
		}
	  }				  
	`)

	req.Var("org_id", org_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["roles"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Role
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
