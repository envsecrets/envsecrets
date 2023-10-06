package subscriptions

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/machinebox/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*Subscription, error)
	GetBySubscriptionID(context.ServiceContext, *clients.GQLClient, string) (*Subscription, error)
	Create(context.ServiceContext, *clients.GQLClient, *CreateOptions) (*Subscription, error)
	List(context.ServiceContext, *clients.GQLClient, *ListOptions) (*[]Subscription, error)
	Update(context.ServiceContext, *clients.GQLClient, string, *UpdateOptions) (*Subscription, error)
	Delete(context.ServiceContext, *clients.GQLClient, string) error
	DeleteBySubscriptionID(context.ServiceContext, *clients.GQLClient, string) error
}

type DefaultService struct{}

// Get a subscription by ID
func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Subscription, error) {

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

	var response struct {
		Subscription Subscription `json:"subscriptions_by_pk"`
	}

	return &response.Subscription, nil
}

// Fetches the row with unique Stripe subscription ID
func (*DefaultService) GetBySubscriptionID(ctx context.ServiceContext, client *clients.GQLClient, subscription_id string) (*Subscription, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp[0], nil
}

// Create a new subscription
func (*DefaultService) Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Subscription, error) {

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

	var response struct {
		Query struct {
			Returning []Subscription `json:"returning"`
		} `json:"insert_subscriptions"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	if len(response.Query.Returning) == 0 {
		return nil, fmt.Errorf("failed to create subscription")
	}

	return &response.Query.Returning[0], nil
}

// List subscriptions
func (*DefaultService) List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) (*[]Subscription, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp []Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update a subscription by ID
func (*DefaultService) Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) (*Subscription, error) {

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
		return nil, err
	}

	//	Unmarshal the response from "returning"
	var resp Subscription
	if err := json.Unmarshal(returning, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Delete a subscription by ID
func (*DefaultService) Delete(ctx context.ServiceContext, client *clients.GQLClient, id string) error {

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
		return errors.New("no rows affected")
	}

	return nil
}

// Delete a subscription by Stripe Subscription ID
func (*DefaultService) DeleteBySubscriptionID(ctx context.ServiceContext, client *clients.GQLClient, id string) error {

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
		return errors.New("no rows affected")
	}

	return nil
}
