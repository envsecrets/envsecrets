package organisations

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/machinebox/graphql"
)

//	Create a new organisation
func Create(ctx context.ServiceContext, client *graphql.Client, options *CreateOptions) (*CreateResponse, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!) {
		insert_organisations(objects: {name: $name}) {
		  returning {
			id
			name
		  }
		}
	  }	  
	`)

	req.Var("name", options.Name)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_organisations"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []CreateResponse
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	result := resp[0]

	//	Add yourself as the first member of the organization
	if _, err = memberships.Create(ctx, client, &memberships.CreateOptions{
		OrgID: result.ID,
	}); err != nil {
		return nil, err
	}

	return &result, nil
}

//	Get a organisation by ID
func Get(ctx context.ServiceContext, client *graphql.Client, id string) (*Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
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

	returning, err := json.Marshal(response["organisations_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	List organisations
func List(ctx context.ServiceContext, client *graphql.Client) (*[]Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery {
		organisations {
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

	returning, err := json.Marshal(response["organisations"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Update a organisation by ID
func Update(ctx context.ServiceContext, client *graphql.Client, id string, options *UpdateOptions) (*Organisation, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $name: String!) {
		update_organisations_by_pk(pk_columns: {id: $id}, _set: {name: $name}) {
			id
		  name
		}
	  }	  
	`)

	req.Var("id", id)
	req.Var("name", options.Name)
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

	var response map[string]interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["update_organisations_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//	Delete a organisation by ID
func Delete(ctx context.ServiceContext, client *graphql.Client, id string) error {
	return nil
}
