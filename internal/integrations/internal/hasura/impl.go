package hasura

import (
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
)

const API = "https://data.pro.hasura.io/v1/graphql"

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Initialize a new GraphQL client.
	client := clients.NewGQLClient2(&clients.GQL2Config{
		BaseURL: API,
		Authorization: &clients.Authorization{
			Token:     options.Credentials["token"].(string),
			TokenType: clients.PAT,
		},
	})

	var query struct {
		Projects []struct {
			Name   string `json:"name"`
			ID     string `json:"id"`
			Tenant struct {
				ID string `json:"id"`
			} `json:"tenant"`
		} `json:"projects"`
	}

	if err := client.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	return &query.Projects, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Initialize a new GraphQL client.
	client := clients.NewGQLClient2(&clients.GQL2Config{
		BaseURL: API,
		Authorization: &clients.Authorization{
			Token:     options.Credentials["token"].(string),
			TokenType: clients.PAT,
		},
	})

	tenant := options.EntityDetails["tenant"].(map[string]interface{})

	//	Get the current hash.
	currentHash, err := getCurrentHash(ctx, client, tenant["id"].(string))
	if err != nil {
		return err
	}

	type UpdateEnvObject struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var envs []UpdateEnvObject
	for key, payload := range *options.Data {
		envs = append(envs, UpdateEnvObject{
			Key:   key,
			Value: payload.Value,
		})
	}

	var mutation struct {
		UpdateTenantEnv struct {
			Hash string
		} `graphql:"updateTenantEnv(tenantId: $tenant_id, currentHash: $hash, envs: $envs)"`
	}

	type uuid string
	if err := client.Mutate(ctx, &mutation, map[string]interface{}{
		"tenant_id": uuid(tenant["id"].(string)),
		"hash":      *currentHash,
		"envs":      envs,
	}); err != nil {
		return err
	}

	if mutation.UpdateTenantEnv.Hash == "" {
		return fmt.Errorf("failed to update tenant variable")
	}

	return nil
}

func getCurrentHash(ctx context.ServiceContext, client *clients.GQLClient2, tenant_id string) (*string, error) {

	var query struct {
		GetTenantEnv struct {
			Hash string
		} `graphql:"getTenantEnv(tenantId: $tenant_id)"`
	}

	type uuid string
	if err := client.Query(ctx, &query, map[string]interface{}{
		"tenant_id": uuid(tenant_id),
	}); err != nil {
		return nil, err
	}

	return &query.GetTenantEnv.Hash, nil
}
