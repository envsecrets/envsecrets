package railway

import (
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/hasura/go-graphql-client"
)

const API = "https://backboard.railway.app/graphql/v2"

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Initialize a new GraphQL client.
	client := clients.NewGQLClient2(&clients.GQL2Config{
		BaseURL: API,
		Authorization: &clients.Authorization{
			Token:     options.Credentials["token"].(string),
			TokenType: clients.Bearer,
		},
	})

	var query struct {
		Me struct {

			//	Get all the projects under the scope of the token.
			Projects struct {
				Edges []struct {
					Node struct {
						Name string
						ID   string

						//	Get all the environments of the project.
						Environments struct {
							Edges []struct {
								Node struct {
									Name string
									ID   string
								}
							}
						}
					}
				}
			}
		}
	}

	err := client.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	//	Transform the response.
	var response []map[string]interface{}
	for _, project := range query.Me.Projects.Edges {
		projectData := map[string]interface{}{
			"name":         project.Node.Name,
			"id":           project.Node.ID,
			"environments": []map[string]interface{}{},
		}

		//	Traverse through all the environments of the project.
		for _, environment := range project.Node.Environments.Edges {
			projectData["environments"] = append(projectData["environments"].([]map[string]interface{}), map[string]interface{}{
				"name": environment.Node.Name,
				"id":   environment.Node.ID,
			})
		}

		response = append(response, projectData)
	}

	return &response, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Initialize a new GraphQL client.
	client := clients.NewGQLClient2(&clients.GQL2Config{
		BaseURL: API,
		Authorization: &clients.Authorization{
			Token:     options.Credentials["token"].(string),
			TokenType: clients.Bearer,
		},
	})

	//type uuid string
	type VariableUpsertInput struct {
		ProjectID     graphql.String `json:"projectId"`
		EnvironmentID graphql.String `json:"environmentId"`
		Name          graphql.String `json:"name"`
		Value         graphql.String `json:"value"`
	}

	project := options.EntityDetails["project"].(map[string]interface{})
	environment := options.EntityDetails["environment"].(map[string]interface{})

	for key, payload := range *options.Data {

		var mutation struct {
			VariableUpsert bool `graphql:"variableUpsert(input: $input)"`
		}

		var inputs = VariableUpsertInput{
			ProjectID:     graphql.String(project["id"].(string)),
			EnvironmentID: graphql.String(environment["id"].(string)),
			Name:          graphql.String(key),
			Value:         graphql.String(payload.Value),
		}

		err := client.Mutate(ctx, &mutation, map[string]interface{}{
			"input": inputs,
		}, graphql.OperationName("MyMutation"))
		if err != nil {
			return err
		}

		if !mutation.VariableUpsert {
			return fmt.Errorf("failed to upsert variable")
		}

	}

	return nil
}
