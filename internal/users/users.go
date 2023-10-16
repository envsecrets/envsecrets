package users

import (
	"encoding/json"
	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*User, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		users(where: {id: {_eq: $id}}) {
			id
			displayName
			email
		}
	  }
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["users"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []User
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

func GetByEmail(ctx context.ServiceContext, client *clients.GQLClient, email string) (*User, error) {

	req := graphql.NewRequest(`
	query MyQuery($email: citext!) {
		users(where: {email: {_eq: $email}}) {
			email
			id
			displayName
		}
	  }	  
	`)

	req.Var("email", email)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["users"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []User
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("no row found")
	}

	return &resp[0], nil
}
