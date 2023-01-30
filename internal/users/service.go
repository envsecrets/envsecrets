package users

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

//	Get a User by ID
func Get(ctx context.Context, client *graphql.Client, id string) (*User, error) {

	var query struct {
		User User `graphql:"user(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.ID(id),
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.User, nil
}

//	Get a User by email
func GetByEmail(ctx context.Context, client *graphql.Client, email string) (*User, error) {

	var query struct {
		User User `graphql:"user(email: $email)"`
	}

	variables := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.User, nil
}
