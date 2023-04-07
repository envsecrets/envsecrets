package users

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/users/commons"
	"github.com/machinebox/graphql"
)

func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.User, *errors.Error) {

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
		return nil, errors.New(err, "failed to marshal json response", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.User
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to nmarshal returning response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

func GetByEmail(ctx context.ServiceContext, client *clients.GQLClient, email string) (*commons.User, *errors.Error) {

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
		return nil, errors.New(err, "failed to marshal json response", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.User
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, "failed to nmarshal returning response", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}
