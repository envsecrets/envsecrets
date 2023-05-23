package vercel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
)

// ---	Flow ---
// 1. Exchange the `code` received from Vercel for an access token: https://api.vercel.com/v2/oauth/access_token
// 2. Save the `access_token` and `installation_id` in Hasura.
func Setup(ctx context.ServiceContext, gqlClient *clients.GQLClient, options *SetupOptions) (*commons.Integration, error) {

	//	Initialize a new HTTP client for Vercel.
	httpClient := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.VercelClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Content-Type",
				Value: "application/x-www-form-urlencoded",
			},
		},
	})

	data := url.Values{}
	data.Set("client_id", os.Getenv("VERCEL_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("VERCEL_CLIENT_SECRET"))
	data.Set("code", options.Code)
	data.Set("redirect_uri", os.Getenv("REDIRECT_DOMAIN")+"/v1/integrations/vercel/setup")

	req, err := http.NewRequest(http.MethodPost, "https://api.vercel.com/v2/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	var response CodeExchangeResponse
	if err := httpClient.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	//	Encrypt the credentials
	credentials, err := commons.EncryptCredentials(ctx, options.OrgID, map[string]interface{}{
		"token_type":   response.TokenType,
		"access_token": response.AccessToken,
		"user_id":      response.UserID,
		"team_id":      response.TeamID,
	})
	if err != nil {
		return nil, err
	}

	//	Create a new record in Hasura.
	return graphql.Insert(ctx, gqlClient, &commons.AddIntegrationOptions{
		OrgID:          options.OrgID,
		InstallationID: response.InstallationID,
		Type:           commons.Vercel,
		Credentials:    base64.StdEncoding.EncodeToString(credentials),
	})
}

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Initialize a new HTTP client for Vercel.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.VercelClientType,
		Authorization: fmt.Sprintf("%v %v", options.Credentials["token_type"], options.Credentials["access_token"]),
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.vercel.com/v9/projects", nil)
	if err != nil {
		return nil, err
	}

	//	If the user had integrated a team account,
	//	then perform ther equest on behalf of that team_id.
	if options.Credentials["team_id"] != nil {
		params := req.URL.Query()
		params.Set("teamId", options.Credentials["team_id"].(string))
		req.URL.RawQuery = params.Encode()
	}

	var response ListProjectsResponse

	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	for index, project := range response.Projects {

		if len(project.LatestDeployments) > 0 {
			project.Username = project.LatestDeployments[0].Creator.Username
			project.LatestDeployments = nil
		}
		response.Projects[index] = project
	}

	return &response.Projects, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Initialize a new HTTP client for Vercel.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.VercelClientType,
		Authorization: fmt.Sprintf("%v %v", options.Credentials["token_type"], options.Credentials["access_token"]),
	})

	//	Initialize TEAM ID
	teamID := options.Credentials["team_id"].(string)

	//	Prepare array of all values
	var array []map[string]interface{}
	for key, value := range options.Secrets {

		//	Prepare the secret type
		var typ string
		v := value.Value

		if value.Type == secretCommons.Ciphertext {

			/* 			//	Create the secret separately in vercel first.
			   			secret, err := CreateSecret(ctx, client, key, value.Value, &teamID)
			   			if err != nil {
			   				return err
			   			}
			*/
			typ = "encrypted"

		} else if value.Type == secretCommons.Plaintext {
			typ = "plain"
		}

		array = append(array, map[string]interface{}{
			"key":    key,
			"value":  v,
			"type":   typ,
			"target": []string{"production", "preview"},
		})
	}

	//	Prepare the request body
	body, err := json.Marshal(array)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.vercel.com/v10/projects/%s/env", options.EntityDetails["id"].(string)), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	//	Prepare Queries
	params := req.URL.Query()
	params.Set("upsert", "true")

	//	If the user had integrated a team account,
	//	then perform ther equest on behalf of that team_id.
	if teamID != "" {
		params.Set("teamId", options.Credentials["team_id"].(string))
	}

	req.URL.RawQuery = params.Encode()

	//	Make the request
	var response VercelResponse

	if err := client.Run(ctx, req, &response); err != nil {
		return err
	}

	if response.Error != nil {
		return fmt.Errorf(response.Error["message"].(string))
	}

	return nil
}

// Creates a new Secret on Vercel.
// Docs: https://vercel.com/docs/rest-api/endpoints#create-a-new-secret
func CreateSecret(ctx context.ServiceContext, client *clients.HTTPClient, name string, value interface{}, teamID *string) (*VercelSecret, error) {

	//	Prepare the request body
	body, err := json.Marshal(map[string]interface{}{
		"name":        name,
		"value":       value,
		"decryptable": true,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.vercel.com/v2/secrets/name", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	//	If the user had integrated a team account,
	//	then perform ther equest on behalf of that team_id.
	if teamID != nil {
		params := req.URL.Query()
		params.Set("teamId", *teamID)
		req.URL.RawQuery = params.Encode()
	}

	//	Make the request
	var response VercelSecret

	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf(response.Error["message"].(string))
	}

	return &response, nil
}
