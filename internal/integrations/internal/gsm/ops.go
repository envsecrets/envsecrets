package gsm

import (
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/envsecrets/envsecrets/internal/context"
	"google.golang.org/api/iterator"
)

func CreateSecret(ctx context.ServiceContext, client *secretmanager.Client, parent string, secretID string) (string, error) {

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create secret: %v", err)
	}
	return result.Name, nil
}

// AddSecretVersion adds a new secret version to the given secret path with the
// provided payload.
func AddSecretVersion(ctx context.ServiceContext, client *secretmanager.Client, path string, payload []byte) (string, error) {

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: path,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	result, err := client.AddSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to add secret version: %v", err)
	}

	return result.Name, nil
}

func ListSecrets(ctx context.ServiceContext, client *secretmanager.Client, parent string) (secrets []string, errors []error) {

	// Build the request.
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}

	// Call the API.
	it := client.ListSecrets(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			errors = append(errors, err)
			secrets = append(secrets, "")
			continue
		}
		secrets = append(secrets, resp.Name)
		errors = append(errors, nil)
	}
	return secrets, errors
}
