package organisations

import (
	"encoding/base64"
	"encoding/json"

	internalErrors "errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/machinebox/graphql"
)

//	Create a new organisation
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Organisation, *errors.Error) {

	errorMessage := "Failed to create organisation"

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

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_organisations"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

//	Create a new organisation
func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Organisation, *errors.Error) {

	errorMessage := "Failed to create organisation"

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $user_id: uuid!) {
		insert_organisations(objects: {name: $name, user_id: $user_id}) {
		  returning {
			id
			name
		  }
		}
	  }
	`)

	req.Var("name", options.Name)
	req.Var("user_id", options.UserID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_organisations"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

//	Get a organisation by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Organisation, *errors.Error) {

	errorMessage := "Failed to fetch the organisation"

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
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

	returning, err := json.Marshal(response["organisations_by_pk"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Get an organisation's key copy of the server
func GetServerKeyCopy(ctx context.ServiceContext, client *clients.GQLClient, id string) ([]byte, *errors.Error) {

	errorMessage := "Failed to fetch server's copy of org's key"

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
			server_copy
		}
	  }	  
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["organisations_by_pk"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	//	Base64 decode the server's key copy
	result, err := base64.StdEncoding.DecodeString(resp.ServerKey)
	if err != nil {
		return nil, errors.New(internalErrors.New(errorMessage), errorMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	return result, nil
}

//	Get a organisation by ID
func GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*Organisation, *errors.Error) {

	errorMessage := "Failed to fetch the organisation"

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		organisations(where: {projects: {environments: {id: {_eq: $env_id}}}}, limit: 1) {
		  id
		}
	  }	   
	`)

	req.Var("env_id", env_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["organisations"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

//	List organisations
func List(ctx context.ServiceContext, client *clients.GQLClient) (*[]Organisation, *errors.Error) {

	errorMessage := "Failed to list organisations"

	req := graphql.NewRequest(`
	query MyQuery {
		organisations {
			id
			name
		}
	  }	  
	`)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["organisations"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Update a organisation by ID
func Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) (*Organisation, *errors.Error) {

	errorMessage := "Failed to update the organisation"

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

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["update_organisations_by_pk"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Organisation
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

func UpdateInviteLimit(ctx context.ServiceContext, client *clients.GQLClient, options *UpdateInviteLimitOptions) *errors.Error {

	errorMessage := "Failed to update the invite limit"

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $limit: Int!) {
		update_organisations(where: {id: {_eq: $id}}, _inc: {invite_limit: $limit}) {
			affected_rows
		  }
	  }			
	`)

	req.Var("id", options.ID)
	req.Var("limit", options.IncrementLimitBy)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["update_organisations"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, errorMessage, errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}

func UpdateServerKeyCopy(ctx context.ServiceContext, client *clients.GQLClient, options *UpdateServerKeyCopyOptions) *errors.Error {

	errorMessage := "Failed to update the invite limit"

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $key: String!) {
		update_organisations(where: {id: {_eq: $id}}, _set: {server_copy: $key}) {
		  affected_rows
		}
	  }				  
	`)

	req.Var("id", options.OrgID)
	req.Var("key", options.Key)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["update_organisations"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(internalErrors.New(errorMessage), errorMessage, errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}

//	Delete a organisation by ID
func Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) error {
	return nil
}
