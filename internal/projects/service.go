package projects

import (
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/hasura/go-graphql-client"
)

//	Create a new workspace
func Create(ctx context.ServiceContext, client *graphql.Client, options *CreateOptions) error {

	var mutation struct {
		name string `graphql:"name(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": &options.Name,
	}

	return client.Mutate(ctx, &mutation, variables)
}

//	Get a workspace by ID
func Get(ctx context.ServiceContext, client *graphql.Client, id string) (*Project, error) {

	var query struct {
		Project Project `graphql:"workspace(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.ID(id),
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.Project, nil
}

//	List projects
func List(ctx context.ServiceContext, client *graphql.Client) (*[]Project, error) {

	var query struct {
		projects []Project `graphql:"projects"`
	}

	if err := client.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	return &query.projects, nil
}

//	Update a workspace by ID
func Update(ctx context.ServiceContext, client *graphql.Client, id string, options *UpdateOptions) error {

	var mutation struct {
		name string `graphql:"name(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": &options.Name,
	}

	return client.Mutate(ctx, &mutation, variables)
}

//	Delete a workspace by ID
func Delete(ctx context.ServiceContext, client *graphql.Client, id string) error {
	return nil
}
