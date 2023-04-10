package subscriptions

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/machinebox/graphql"
)

//	Create a new workspace
func Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Subscription, *errors.Error) {

	errorMessage := "Failed to create subscription"

	req := graphql.NewRequest(`
	mutation MyMutation($subscription_id: String!, $org_id: uuid!) {
		insert_subscriptions(objects: {subscription_id: $subscription_id, org_id: $org_id}) {
		  returning {
			id
		  }
		}
	  }	  
	`)

	req.Var("subscription_id", options.SubscriptionID)
	req.Var("org_id", options.OrgID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["insert_subscriptions"].(map[string]interface{})["returning"].([]interface{}))
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

//	Get a workspace by ID
func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Subscription, *errors.Error) {

	errorMessage := "Failed to fetch subscription"

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		subscriptions_by_pk(id: $id) {
			id
			org_id
			status
			subscription_id
		}
	  }	  
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["subscriptions_by_pk"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

func GetBySubscriptionID(ctx context.ServiceContext, client *clients.GQLClient, subscription_id string) (*Subscription, *errors.Error) {

	errorMessage := "Failed to fetch subscription"

	req := graphql.NewRequest(`
	query MyQuery($subscription_id: String!) {
		subscriptions(where: {subscription_id: {_eq: $subscription_id}}) {
			id
			org_id
			status
		  }
	  }	  
	`)

	req.Var("subscription_id", subscription_id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["subscriptions"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp[0], nil
}

//	List subscriptions
func List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) (*[]Subscription, *errors.Error) {

	errorMessage := "Failed to list subscriptions"

	req := graphql.NewRequest(`
	query MyQuery($org_id: uuid!) {
		subscriptions(where: {org_id: {_eq: $org_id}}) {
		  status
		}
	  }	  
	`)

	req.Var("org_id", options.OrgID)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["subscriptions"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp []Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Update a workspace by ID
func Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) (*Subscription, *errors.Error) {

	errorMessage := "Failed to update the subscription"

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $status: String!) {
		update_subscriptions_by_pk(pk_columns: {id: $id}, _set: {status: $status}) {
			id
		  status
		}
	  }	  
	`)

	req.Var("id", id)
	req.Var("status", options.Status)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	returning, err := json.Marshal(response["update_subscriptions_by_pk"])
	if err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Unmarshal the response from "returning"
	var resp Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, errors.New(err, errorMessage, errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &resp, nil
}

//	Delete a subscription by ID
func Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) *errors.Error {

	errorMessage := "Failed to delete the subscription"

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!) {
		delete_subscriptions(where: {id: {_eq: $id}}) {
		  affected_rows
		}
	  }			
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["delete_subscriptions"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, errorMessage, errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}

//	Delete a subscription by Stripe Subscription ID
func DeleteBySubscriptionID(ctx context.ServiceContext, client *clients.GQLClient, id string) *errors.Error {

	errorMessage := "Failed to delete the subscription"

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!) {
		delete_subscriptions(where: {subscription_id: {_eq: $id}}) {
		  affected_rows
		}
	  }			
	`)

	req.Var("id", id)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["delete_subscriptions"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, errorMessage, errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}
