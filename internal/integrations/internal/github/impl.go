package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
)

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Get installation's access token
	auth, err := GetInstallationAccessToken(ctx, options.InstallationID)
	if err != nil {
		return nil, err
	}

	//	Initialize a new HTTP client for Github.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.GithubClientType,
		Authorization: "Bearer " + auth.Token,
	})

	//	Fetch the repositories
	respositoryResponse, err := ListRepositories(ctx, client)
	if err != nil {
		return nil, err
	}

	return &respositoryResponse.Repositories, nil
}

func ListRepositories(ctx context.ServiceContext, client *clients.HTTPClient) (*ListRepositoriesResponse, error) {

	//	Get user's access token from Github API.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/installation/repositories", nil)
	if err != nil {
		return nil, err
	}

	var repositoriesResponse ListRepositoriesResponse
	err = client.Run(ctx, req, &repositoriesResponse)
	if err != nil {
		return nil, err
	}

	return &repositoriesResponse, nil
}

// -- Flow --
// 1. Get repository's action secrets public key.
// 2. Encrypt the secret data.
// 3. Post the secrets to Github actions endpoint.
func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Get installation's access token
	auth, err := GetInstallationAccessToken(ctx, options.InstallationID)
	if err != nil {
		return err
	}

	//	Initialize a new HTTP client for Github.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.GithubClientType,
		Authorization: "Bearer " + auth.Token,
	})

	//	Extract the slug from entity details
	slug := options.EntityDetails["full_name"].(string)

	for key, payload := range *options.Data {

		//	If the payload is of type `ciphertext`,
		//	we have to encrypt its value and push it to Github action's secrets.
		if !payload.IsExposable() {

			//	Get the public key.
			publicKey, err := getRepositoryActionsSecretsPublicKey(ctx, client, slug)
			if err != nil {
				return err
			}

			//	Base64 decode the value of the secrets.
			payload.Decode()

			//	Encrypt the secret value.
			encryptedValue, err := encryptSecret(publicKey.Key, payload.GetValue())
			if err != nil {
				return err
			}

			//	Add response handler to HTTP client.
			client.ResponseHandler = func(response *http.Response) error {

				//	Github Responses:
				//	201 -> New secret created
				//	204 -> Existing secret updated
				if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
					return errors.New("failed to push secret to github repo")
				}
				return nil
			}

			//	Post the secret to Github actions.
			if err := pushRepositorySecret(ctx, client, slug, key, publicKey.KeyID, encryptedValue); err != nil {
				return err
			}

		} else {

			//	Add response handler to HTTP client.
			client.ResponseHandler = func(response *http.Response) error {

				//	Github Responses:
				//	201 (Created) -> New variable created
				//	409 (Conflict) -> Variable exists
				if response.StatusCode == http.StatusConflict {

					//	Delete the variable and recreate it.
					if err := deleteRepositoryVariable(ctx, client, slug, key); err != nil {
						return err
					}

					return pushRepositoryVariable(ctx, client, slug, key, payload.Value)
				}

				return nil
			}

			//	If the payload type is `plaintext`,
			//	save it as a normal variable in Github actions.
			if err := pushRepositoryVariable(ctx, client, slug, key, payload.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

func pushRepositorySecret(ctx context.ServiceContext, client *clients.HTTPClient, slug, secretName, keyID, value string) error {

	body, err := json.Marshal(map[string]interface{}{
		"encrypted_value": value,
		"key_id":          keyID,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("https://api.github.com/repos/%s/actions/secrets/%s", slug, secretName), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return client.Run(ctx, req, nil)
}

func pushRepositoryVariable(ctx context.ServiceContext, client *clients.HTTPClient, slug, name, value string) error {

	body, err := json.Marshal(map[string]interface{}{
		"name":  name,
		"value": value,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.github.com/repos/%s/actions/variables", slug), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return client.Run(ctx, req, nil)
}

func deleteRepositoryVariable(ctx context.ServiceContext, client *clients.HTTPClient, slug, name string) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("https://api.github.com/repos/%s/actions/variables/%s", slug, name), nil)
	if err != nil {
		return err
	}

	return client.Run(ctx, req, nil)
}

func GetInstallationAccessToken(ctx context.ServiceContext, installationID string) (*InstallationAccessTokenResponse, error) {

	//	Load the github private key
	key := os.Getenv("GITHUB_PRIVATE_KEY")

	//	Base64 decode the key to get PEM data
	value, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	//	Authenticate as a github app
	jwt, err := generateGithuAppJWT(value)
	if err != nil {
		return nil, err
	}

	//	Initialize a new HTTP client for Github.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.GithubClientType,
		Authorization: "Bearer " + jwt,
	})

	//	Get user's access token from Github API.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID), nil)
	if err != nil {
		return nil, err
	}

	var response InstallationAccessTokenResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Fetches the public key for action secrets for supplied repository slug.
func getRepositoryActionsSecretsPublicKey(ctx context.ServiceContext, client *clients.HTTPClient, slug string) (*RepositoryActionsSecretsPublicKeyResponse, error) {

	//	Get user's access token from Github API.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/actions/secrets/public-key", slug), nil)
	if err != nil {
		return nil, err
	}

	var response RepositoryActionsSecretsPublicKeyResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
