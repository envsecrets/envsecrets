package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
)

func Setup(ctx context.ServiceContext, client *clients.GQLClient, options *SetupOptions) *errors.Error {

	//	Create a new record in Hasura.
	if err := graphql.Insert(ctx, client, &commons.AddIntegrationOptions{
		OrgID:          options.OrgID,
		InstallationID: options.InstallationID,
		Type:           commons.Github,
	}); err != nil {
		return err
	}

	//	TODO: Redirect the user to front-end to complete post-integration steps.
	return nil
}

func ListEntities(ctx context.ServiceContext, integration *commons.Integration) (*commons.Entities, *errors.Error) {

	//	Get installation's access token
	auth, err := getInstallationAccessToken(ctx, integration.InstallationID)
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

	//	Convert to entities_response and return
	var entities commons.Entities
	for _, item := range respositoryResponse.Repositories {
		entity := *item.ToEntity()
		entity.InstallationID = integration.InstallationID
		entities = append(entities, entity)
	}

	return &entities, nil
}

func ListRepositories(ctx context.ServiceContext, client *clients.HTTPClient) (*ListRepositoriesResponse, *errors.Error) {

	//	Get user's access token from Github API.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/installation/repositories", nil)
	if err != nil {
		return nil, errors.New(err, "failed prepare oauth access token request", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	repositoriesResponsePayload, er := client.Run(ctx, req)
	if er != nil {
		return nil, er
	}

	defer repositoriesResponsePayload.Body.Close()

	repositoriesResponseBody, err := ioutil.ReadAll(repositoriesResponsePayload.Body)
	if err != nil {
		return nil, errors.New(err, "failed to read github repositories response body", errors.ErrorTypeBadResponse, errors.ErrorSourceGithub)
	}

	var repositoriesResponse ListRepositoriesResponse
	if err := json.Unmarshal(repositoriesResponseBody, &repositoriesResponse); err != nil {
		return nil, errors.New(err, "failed to unmarshal github repositories response body", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &repositoriesResponse, nil
}

func getInstallationAccessToken(ctx context.ServiceContext, installationID string) (*InstallationAccessTokenResponse, *errors.Error) {

	//	Authenticate as a github app
	jwt, err := generateGithuAppJWT("keys/github-private-key.pem")
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

	authResponsePayload, err := client.Run(ctx, req)
	if err != nil {
		return nil, err
	}

	defer authResponsePayload.Body.Close()

	authResponseBody, er := ioutil.ReadAll(authResponsePayload.Body)
	if er != nil {
		return nil, errors.New(er, "failed to read github access token response body", errors.ErrorTypeBadResponse, errors.ErrorSourceGithub)
	}

	var authResponse InstallationAccessTokenResponse
	if err := json.Unmarshal(authResponseBody, &authResponse); err != nil {
		return nil, errors.New(err, "failed to unmarshal github access token response body", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	return &authResponse, nil
}
