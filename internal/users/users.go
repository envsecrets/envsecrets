package users

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/users/commons"
	"github.com/machinebox/graphql"
)

func Get(ctx context.ServiceContext, client *client.GQLClient, id string) (*commons.User, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		users_by_pk(id: $id) {
			id
			name
			email
		}
	  }
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err.Error
	}

	returning, err := json.Marshal(response["users_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp commons.User
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetByEmail(ctx context.ServiceContext, client *client.GQLClient, email string) (*commons.User, error) {

	req := graphql.NewRequest(`
	query MyQuery($email: String!) {
		users_by_pk(email: $email) {
			id
			name
			email
		}
	  }	  
	`)

	req.Var("email", email)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err.Error
	}

	returning, err := json.Marshal(response["users_by_pk"])
	if err != nil {
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp commons.User
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
