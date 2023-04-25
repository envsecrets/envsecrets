package graphql

import (
	"encoding/base64"
	"encoding/json"

	internalErrors "errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/machinebox/graphql"
)

//	Create a new key
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateOptions) *errors.Error {

	errorMessage := "Failed to create key"

	req := graphql.NewRequest(`
	mutation MyMutation($public_key: String!, $private_key: String!, $protected_key: String!, $salt: String!) {
		insert_keys(objects: {private_key: $private_key, protected_key: $protected_key, public_key: $public_key, salt: $salt}) {
		  affected_rows
		}
	  }	  
	`)

	req.Var("public_key", options.PublicKey)
	req.Var("private_key", options.PrivateKey)
	req.Var("protected_key", options.ProtectedKey)
	req.Var("salt", options.Salt)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_keys"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, errorMessage, errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}

//	Create a new key with User ID
func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateWithUserIDOptions) *errors.Error {

	errorMessage := "Failed to create key"

	req := graphql.NewRequest(`
	mutation MyMutation($public_key: String!, $private_key: String!, $protected_key: String!, $salt: String!, $user_id: uuid!) {
		insert_keys(objects: {private_key: $private_key, protected_key: $protected_key, public_key: $public_key, salt: $salt, user_id: $user_id}) {
		  affected_rows
		}
	  }			
	`)

	req.Var("public_key", options.PublicKey)
	req.Var("private_key", options.PrivateKey)
	req.Var("protected_key", options.ProtectedKey)
	req.Var("salt", options.Salt)
	req.Var("user_id", options.UserID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_keys"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, errorMessage, errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}

//	Get a key by User ID
func GetByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) (*commons.Key, *errors.Error) {

	errorMessage := "Failed to fetch the key"

	req := graphql.NewRequest(`
	query MyQuery($user_id: uuid!) {
		keys(where: {user_id: {_eq: $user_id}}) {
		  private_key
		  protected_key
		  public_key
		  salt
		  id
		}
	  }			
	`)

	req.Var("user_id", user_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["keys"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Key
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	if len(resp) == 0 {
		return nil, errors.New(internalErrors.New(errorMessage), errorMessage, errors.ErrorTypeDoesNotExist, errors.ErrorSourceGraphQL)
	}

	return &resp[0], nil
}

//	Fetches only the public key by User ID
func GetPublicKeyByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) ([]byte, *errors.Error) {

	errorMessage := "Failed to fetch the key"

	req := graphql.NewRequest(`
	query MyQuery($user_id: uuid!) {
		keys(where: {user_id: {_eq: $user_id}}) {
		  public_key
		}
	  }			
	`)

	req.Var("user_id", user_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["keys"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Key
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	result, err := base64.StdEncoding.DecodeString(resp[0].PublicKey)
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeBadResponse, errors.ErrorSourceGo)
	}

	return result, nil
}

//	Fetches only the public key by User's email
func GetPublicKeyByUserEmail(ctx context.ServiceContext, client *clients.GQLClient, email string) ([]byte, *errors.Error) {

	errorMessage := "Failed to fetch the key"

	req := graphql.NewRequest(`
	query MyQuery($email: citext!) {
		keys(where: {user: {email: {_eq: $email}}}) {
		  public_key
		}
	  }				  
	`)

	req.Var("email", email)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["keys"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Key
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	result, err := base64.StdEncoding.DecodeString(resp[0].PublicKey)
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeBadResponse, errors.ErrorSourceGo)
	}

	return result, nil
}
