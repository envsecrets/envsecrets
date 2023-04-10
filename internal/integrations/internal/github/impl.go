package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
)

func Setup(ctx context.ServiceContext, client *clients.GQLClient, options *SetupOptions) (*commons.Integration, *errors.Error) {

	//	Create a new record in Hasura.
	return graphql.Insert(ctx, client, &commons.AddIntegrationOptions{
		OrgID:          options.OrgID,
		InstallationID: options.InstallationID,
		Type:           commons.Github,
	})
}

func ListEntities(ctx context.ServiceContext, integration *commons.Integration) (interface{}, *errors.Error) {

	//	Get installation's access token
	auth, err := GetInstallationAccessToken(ctx, integration.InstallationID)
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

func ListRepositories(ctx context.ServiceContext, client *clients.HTTPClient) (*ListRepositoriesResponse, *errors.Error) {

	//	Get user's access token from Github API.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/installation/repositories", nil)
	if err != nil {
		return nil, errors.New(err, "failed prepare oauth access token request", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	var repositoriesResponse ListRepositoriesResponse
	er := client.Run(ctx, req, &repositoriesResponse)
	if er != nil {
		return nil, er
	}

	return &repositoriesResponse, nil
}

//	-- Flow --
//	1. Get repository's action secrets public key.
//	2. Encrypt the secret data.
//	3. Post the secrets to Github actions endpoint.
func Sync(ctx context.ServiceContext, options *commons.SyncOptions) *errors.Error {

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

	for key, payload := range options.Data {

		//	If the payload is of type `ciphertext`,
		//	we have to encrypt its value and push it to Github action's secrets.
		if payload.Type == secretCommons.Ciphertext {

			//	Get the public key.
			publicKey, err := getRepositoryActionsSecretsPublicKey(ctx, client, slug)
			if err != nil {
				return err
			}

			//	Encrypt the secret value.
			encryptedValue, er := encryptSecret(publicKey.Key, payload.Value.(string))
			if er != nil {
				return errors.New(er, "failed to encrypt secret", errors.ErrorTypeBadResponse, errors.ErrorSourceGo)
			}

			//	Add response handler to HTTP client.
			client.ResponseHandler = func(response *http.Response) *errors.Error {

				//	Github Responses:
				//	201 -> New secret created
				//	204 -> Existing secret updated
				if response.StatusCode != 201 && response.StatusCode != 204 {
					return errors.New(fmt.Errorf(fmt.Sprint(response.StatusCode)), "failed to push secret to github repo", errors.ErrorTypeBadResponse, errors.ErrorSourceGithub)
				}
				return nil
			}

			//	Post the secret to Github actions.
			if err := pushRepositorySecret(ctx, client, slug, key, publicKey.KeyID, encryptedValue); err != nil {
				return err
			}

		} else if payload.Type == secretCommons.Plaintext {

			//	Add response handler to HTTP client.
			client.ResponseHandler = func(response *http.Response) *errors.Error {

				//	Github Responses:
				//	201 (Created) -> New variable created
				//	409 (Conflict) -> Variable exists
				if response.StatusCode == 409 {

					//	Delete the variable and recreate it.
					if err := deleteRepositoryVariable(ctx, client, slug, key); err != nil {
						return err
					}

					return pushRepositoryVariable(ctx, client, slug, key, payload.Value.(string))
				}

				return nil
			}

			//	If the payload type is `plaintext`,
			//	save it as a normal variable in Github actions.
			if err := pushRepositoryVariable(ctx, client, slug, key, payload.Value.(string)); err != nil {
				return err
			}
		}
	}

	return nil
}

func pushRepositorySecret(ctx context.ServiceContext, client *clients.HTTPClient, slug, secretName, keyID, value string) *errors.Error {

	body, err := json.Marshal(map[string]interface{}{
		"encrypted_value": value,
		"key_id":          keyID,
	})

	if err != nil {
		return errors.New(err, "failed prepare json body to push secrets to github repo", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("https://api.github.com/repos/%s/actions/secrets/%s", slug, secretName), bytes.NewBuffer(body))
	if err != nil {
		return errors.New(err, "failed prepare http request to push secrets to github repo", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return client.Run(ctx, req, nil)
}

func pushRepositoryVariable(ctx context.ServiceContext, client *clients.HTTPClient, slug, name, value string) *errors.Error {

	body, err := json.Marshal(map[string]interface{}{
		"name":  name,
		"value": value,
	})

	if err != nil {
		return errors.New(err, "failed prepare json body to push variables to github repo", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.github.com/repos/%s/actions/variables", slug), bytes.NewBuffer(body))
	if err != nil {
		return errors.New(err, "failed prepare http request to push variables to github repo", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return client.Run(ctx, req, nil)
}

func deleteRepositoryVariable(ctx context.ServiceContext, client *clients.HTTPClient, slug, name string) *errors.Error {

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("https://api.github.com/repos/%s/actions/variables/%s", slug, name), nil)
	if err != nil {
		return errors.New(err, "failed prepare http request to push variables to github repo", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return client.Run(ctx, req, nil)
}

func GetInstallationAccessToken(ctx context.ServiceContext, installationID string) (*InstallationAccessTokenResponse, *errors.Error) {

	//	Load the github private key
	key := os.Getenv("GITHUB_PRIVATE_KEY")

	//	Base64 decode the key to get PEM data
	value, er := base64.StdEncoding.DecodeString(key)
	if er != nil {
		return nil, errors.New(er, "failed to base64 decode PEM file", errors.ErrorTypeDoesNotExist, errors.ErrorSourceGo)
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
	req, er := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID), nil)
	if err != nil {
		return nil, errors.New(er, "failed prepare oauth access token request", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	var response InstallationAccessTokenResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

//	Fetches the public key for action secrets for supplied repository slug.
func getRepositoryActionsSecretsPublicKey(ctx context.ServiceContext, client *clients.HTTPClient, slug string) (*RepositoryActionsSecretsPublicKeyResponse, *errors.Error) {

	//	Get user's access token from Github API.
	req, er := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/actions/secrets/public-key", slug), nil)
	if er != nil {
		return nil, errors.New(er, "failed prepare repository actions secret public key request", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	var response RepositoryActionsSecretsPublicKeyResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}