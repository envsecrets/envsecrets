package vercel

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
)

//	---	Flow ---
//	1. Exchange the `code` received from Vercel for an access token: https://api.vercel.com/v2/oauth/access_token
//	2. Save the `access_token` and `installation_id` in Hasura.
func Setup(ctx context.ServiceContext, gqlClient *clients.GQLClient, options *SetupOptions) *errors.Error {

	//	Initialize a new HTTP client for Vercel.
	httpClient := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.HTTPClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Content-Type",
				Value: "application/x-www-form-urlencoded",
			},
		},
	})

	body, err := json.Marshal(map[string]interface{}{
		"client_id":     os.Getenv("VERCEL_CLIENT_ID"),
		"client_secret": os.Getenv("VERCEL_CLIENT_SECRET"),
		"code":          options.Code,
		"redirect_uri":  "",
	})
	if err != nil {
		return errors.New(err, "failed to marshal request body", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.vercel.com/v2/oauth/access_token", bytes.NewBuffer(body))
	if err != nil {
		return errors.New(err, "failed to prepare http request", errors.ErrorTypeRequestFailed, errors.ErrorSourceGo)
	}

	var response CodeExchangeResponse
	if err := httpClient.Run(ctx, req, &response); err != nil {
		return err
	}

	//	Create a new record in Hasura.
	if err := graphql.Insert(ctx, gqlClient, &commons.AddIntegrationOptions{
		OrgID:          options.OrgID,
		InstallationID: response.InstallationID,
		Type:           commons.Vercel,
	}); err != nil {
		return err
	}

	return nil
}

func ListEntities(ctx context.ServiceContext, integration *commons.Integration) (interface{}, *errors.Error) {

	//	Initialize a new HTTP client for Github.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.GithubClientType,
		//	Authorization: "Bearer " + auth.Token,
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

	return nil
}
