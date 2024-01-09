package graphql

import (
	"encoding/base64"
	"encoding/json"

	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/machinebox/graphql"
)

// Create a new key
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($public_key: String!, $private_key: String!, $protected_key: String!, $sync_key: String!, $salt: String!) {
		insert_keys(objects: {private_key: $private_key, protected_key: $protected_key, public_key: $public_key, salt: $salt, sync_key: $sync_key}) {
		  affected_rows
		}
	  }	  
	`)

	req.Var("public_key", options.PublicKey)
	req.Var("private_key", options.PrivateKey)
	req.Var("protected_key", options.ProtectedKey)
	req.Var("sync_key", options.SyncKey)
	req.Var("salt", options.Salt)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_keys"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to create key")
	}

	return nil
}

// Create a new key with User ID
func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateWithUserIDOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($public_key: String!, $private_key: String!, $protected_key: String!, $sync_key: String!, $salt: String!, $user_id: uuid!) {
		insert_keys(objects: {private_key: $private_key, protected_key: $protected_key, public_key: $public_key, salt: $salt, user_id: $user_id, sync_key: $sync_key}) {
		  affected_rows
		}
	  }				  
	`)

	req.Var("public_key", options.PublicKey)
	req.Var("private_key", options.PrivateKey)
	req.Var("protected_key", options.ProtectedKey)
	req.Var("sync_key", options.SyncKey)
	req.Var("salt", options.Salt)
	req.Var("user_id", options.UserID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_keys"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to create key")
	}

	return nil
}

// Create a new key with User ID
func CreateSyncKey(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateSyncKeyOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($sync_key: String!) {
		insert_keys(objects: {sync_key: $sync_key}) {
		  affected_rows
		}
	  }				  
	`)

	req.Var("sync_key", options.SyncKey)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_keys"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to create key")
	}

	return nil
}

// Get a key by User ID
func GetByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) (*commons.Key, error) {

	req := graphql.NewRequest(`
	query MyQuery($user_id: uuid!) {
		keys(where: {user_id: {_eq: $user_id}}) {
		  private_key
		  protected_key
		  public_key
		  sync_key
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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Key
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("failed to fetch the key")

	}

	return &resp[0], nil
}

// Fetches only the public key by User ID
func GetPublicKeyByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) ([]byte, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Key
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("failed to fetch the key")
	}

	result, err := base64.StdEncoding.DecodeString(resp[0].PublicKey)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Fetches only the public key by User's email
func GetPublicKeyByUserEmail(ctx context.ServiceContext, client *clients.GQLClient, email string) ([]byte, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Key
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("failed to fetch the key")
	}

	result, err := base64.StdEncoding.DecodeString(resp[0].PublicKey)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Fetches only the sync key by User ID
func GetSyncKeyByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) ([]byte, error) {

	req := graphql.NewRequest(`
	query MyQuery($user_id: uuid!) {
		keys(where: {user_id: {_eq: $user_id}}) {
		  sync_key
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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []commons.Key
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("failed to fetch the key")
	}

	result, err := base64.StdEncoding.DecodeString(resp[0].SyncKey)
	if err != nil {
		return nil, err
	}

	return result, nil
}
