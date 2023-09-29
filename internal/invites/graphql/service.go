package graphql

import (
	"encoding/json"
	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/machinebox/graphql"
)

// Get a invite by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Invite, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		invites_by_pk(id: $id) {
			id
			key
			org_id
			role_id
			email
			accepted
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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp commons.Invite
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func Insert(ctx context.ServiceContext, client *clients.GQLClient, objects []commons.InsertOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($objects: [invites_insert_input!]!) {
		insert_invites(objects: $objects) {
		  affected_rows
		}
	  }			
	`)

	req.Var("objects", objects)

	var response struct {
		Query struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"insert_invites"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	//	Validate the mutation as been written to the database
	if response.Query.AffectedRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

// List invites
func List(ctx context.ServiceContext, client *clients.GQLClient, options *commons.ListOptions) (*[]commons.Invite, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Invite
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update a invite by ID
func Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *commons.UpdateOptions) (*commons.Invite, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $accepted: Boolean!) {
		update_invites(where: {id: {_eq: $id}}, _set: {accepted: $accepted}) {
		  returning {
			id
			accepted
		  }
		}
	  }	   
	`)

	req.Var("id", id)
	req.Var("accepted", options.Accepted)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning := response["update_invites"].(map[string]interface{})["returning"].([]interface{})

	data, err := json.Marshal(returning[0])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp commons.Invite
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
