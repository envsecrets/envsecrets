package secrets

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

// Write a secret
func Put(ctx context.Context, client *vault.Client, path string, data map[string]interface{}) error {

	/* 	//	Get existing secret
	   	secret, err := Get(ctx, client, path, nil)
	   	if err != nil {
	   		return err
	   	}

	   	payload := secret.Data
	   	for key, value := range data {
	   		payload[key] = value
	   	}
	*/
	_, err := client.KVv2("secret").Patch(ctx, path, data)
	return err
}

// Get a secret
func Get(ctx context.Context, client *vault.Client, path string, version *int) (*vault.KVSecret, error) {

	if version != nil {
		client.KVv2("secret").GetVersion(ctx, path, *version)
	}

	return client.KVv2("secret").Get(ctx, path)
}

// Get all versions of a secret
func GetVersions(ctx context.Context, client *vault.Client, path string) ([]vault.KVVersionMetadata, error) {
	return client.KVv2("secret").GetVersionsAsList(ctx, path)
}

// List secret
func List(ctx context.Context, client *vault.Client, path string) (*vault.Secret, error) {
	return client.Logical().List("secret/metadata/" + path)
}
