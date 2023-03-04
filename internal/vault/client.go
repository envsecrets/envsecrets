package vault

import (
	"os"

	"github.com/hashicorp/vault/api"
)

var DClient *api.Client

func init() {

	DClient, _ = NewClient()
}

func NewClient() (*api.Client, error) {

	config := api.DefaultConfig() // modify for more granular configuration
	config.Address = os.Getenv("VAULT_ADDRESS")
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Authenticate
	client.SetToken(os.Getenv("VAULT_ROOT_TOKEN"))

	return client, nil
}
