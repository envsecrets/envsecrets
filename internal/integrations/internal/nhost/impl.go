package nhost

import (
	"fmt"

	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/hasura/go-graphql-client"
)

const GRAPHQL_BASE_URL = "https://nhost.run/v1/graphql"
const AUTH_BASE_URL = "https://nhost.run/v1/auth"

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Initialize a new GraphQL client.
	client, err := getClient(ctx, options.Credentials["token"].(string))
	if err != nil {
		return nil, err
	}

	var query struct {
		Apps []struct {
			ID        string `json:"id"`
			Slug      string `json:"slug"`
			Workspace struct {
				Slug string `json:"slug"`
			} `json:"workspace"`
		} `json:"apps"`
	}

	if err := client.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	//	Transform the response.
	var response []map[string]interface{}
	for _, app := range query.Apps {
		response = append(response, map[string]interface{}{
			"id":   app.ID,
			"name": fmt.Sprintf("%s/%s", app.Workspace.Slug, app.Slug),
		})
	}

	return &response, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Initialize a new GraphQL client.
	client, err := getClient(ctx, options.Credentials["token"].(string))
	if err != nil {
		return err
	}

	type ConfigEnvironmentVariable struct {
		Name  graphql.String `json:"name" graphql:"name"`
		Value graphql.String `json:"value" graphql:"value"`
	}
	type ConfigEnvironmentVariableInsertInput ConfigEnvironmentVariable

	appID := options.EntityDetails["id"].(string)

	for key, payload := range *options.Data {

		var mutation struct {
			InsertSecret ConfigEnvironmentVariable `graphql:"insertSecret(appID: $app_id, secret: $secret)"`
		}

		type uuid string
		if err := client.Mutate(ctx, &mutation, map[string]interface{}{
			"app_id": uuid(appID),
			"secret": ConfigEnvironmentVariableInsertInput{
				Name:  graphql.String(key),
				Value: graphql.String(payload.Value),
			},
		}); err != nil {

			//	One possibility of error is that the secret already exists,
			//	in which case we need to attempt to update the value.

			var mutation struct {
				UpdateSecret ConfigEnvironmentVariable `graphql:"updateSecret(appID: $app_id, secret: $secret)"`
			}

			if err := client.Mutate(ctx, &mutation, map[string]interface{}{
				"app_id": uuid(appID),
				"secret": ConfigEnvironmentVariableInsertInput{
					Name:  graphql.String(key),
					Value: graphql.String(payload.Value),
				},
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func getClient(ctx context.ServiceContext, token string) (*clients.GQLClient2, error) {

	//	Exchange the PAT for a JWT.
	session, err := auth.GetService().SigninWithPAT(ctx, clients.NewNhostClient(&clients.NhostConfig{
		BaseURL: AUTH_BASE_URL,
	}), &auth.SigninWithPATOptions{
		PAT: token,
	})
	if err != nil {
		return nil, err
	}

	//	Initialize a new GraphQL client.
	client := clients.NewGQLClient2(&clients.GQL2Config{
		BaseURL: GRAPHQL_BASE_URL,
		Authorization: &clients.Authorization{
			Token:     session.Session["accessToken"].(string),
			TokenType: clients.Bearer,
		},
	})

	return client, nil
}
