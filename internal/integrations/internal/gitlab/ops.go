package gitlab

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
)

// Fetches the list of projects from Gitlab.
func ListProjects(ctx context.ServiceContext, client *clients.HTTPClient) (*ListProjectsResponse, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/api/v4/projects?membership=true&simple=true", nil)
	if err != nil {
		return nil, err
	}

	var response ListProjectsResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Fetches the list of groups from Gitlab.
func ListGroups(ctx context.ServiceContext, client *clients.HTTPClient) (*ListGroupsResponse, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/api/v4/groups?owned=true", nil)
	if err != nil {
		return nil, err
	}

	var response ListGroupsResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Creates a new project variable.
func CreateProjectVariable(ctx context.ServiceContext, client *clients.HTTPClient, options *CreateVariableOptions) (*Variable, error) {

	URL := fmt.Sprintf("https://gitlab.com/api/v4/projects/%v/variables", options.ID)

	body, err := options.Variable.Marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var response Variable
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	//	If the returned payload is nil,
	//	then just assume that the variable already exists,
	//	and try to update the value.
	if response.Key == "" {
		return UpdateProjectVariable(ctx, client, options)
	}

	return &response, nil
}

// Creates a new group variable.
func CreateGroupVariable(ctx context.ServiceContext, client *clients.HTTPClient, options *CreateVariableOptions) (*Variable, error) {

	URL := fmt.Sprintf("https://gitlab.com/api/v4/groups/%v/variables", options.ID)

	body, err := options.Variable.Marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var response Variable
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	//	If the returned payload is nil,
	//	then just assume that the variable already exists,
	//	and try to update the value.
	if response.Key == "" {
		return UpdateGroupVariable(ctx, client, options)
	}

	return &response, nil
}

// Updates an existing variable.
func UpdateProjectVariable(ctx context.ServiceContext, client *clients.HTTPClient, options *CreateVariableOptions) (*Variable, error) {

	URL := fmt.Sprintf("https://gitlab.com/api/v4/projects/%v/variables/%s", options.ID, options.Variable.Key)

	body, err := options.Variable.Marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var response Variable
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Updates an existing variable.
func UpdateGroupVariable(ctx context.ServiceContext, client *clients.HTTPClient, options *CreateVariableOptions) (*Variable, error) {

	URL := fmt.Sprintf("https://gitlab.com/api/v4/groups/%v/variables/%s", options.ID, options.Variable.Key)

	body, err := options.Variable.Marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var response Variable
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetAccessToken(ctx context.ServiceContext, options *TokenRequestOptions) (*TokenResponse, error) {

	//	Initialize a new HTTP client.
	httpClient := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.HTTPClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Content-Type",
				Value: "application/json",
			},
		},
	})

	payload := map[string]interface{}{
		"client_id":     os.Getenv("GITLAB_APP_ID"),
		"client_secret": os.Getenv("GITLAB_APP_SECRET"),
	}

	if options.Code != "" {
		payload["code"] = options.Code
		payload["grant_type"] = "authorization_code"
	} else if options.RefreshToken != "" {
		payload["refresh_token"] = options.RefreshToken
		payload["grant_type"] = "refresh_token"
	} else {
		return nil, errors.New("either code or refresh token is required")
	}

	if options.RedirectURI != "" {
		payload["redirect_uri"] = options.RedirectURI
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://gitlab.com/oauth/token", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var response TokenResponse
	if err := httpClient.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func RefreshToken(ctx context.ServiceContext, options *TokenRefreshOptions) (*TokenResponse, error) {

	//	Generate a fresh pair of tokens
	tokens, err := GetAccessToken(ctx, &TokenRequestOptions{
		RefreshToken: options.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	//	Save updated credentials in Hasura.
	if tokens.RefreshToken != "" {

		//	Encrypt the credentials
		credentials, err := commons.EncryptCredentials(ctx, options.OrgID, map[string]interface{}{
			"token_type":    tokens.TokenType,
			"refresh_token": tokens.RefreshToken,
		})
		if err != nil {
			return nil, err
		}

		//	Initialize Hasura client with admin privileges
		client := clients.NewGQLClient(&clients.GQLConfig{
			Type: clients.HasuraClientType,
			Headers: []clients.Header{
				clients.XHasuraAdminSecretHeader,
			},
		})

		err = graphql.UpdateCredentials(ctx, client, &commons.UpdateCredentialsOptions{
			ID:          options.IntegrationID,
			Credentials: base64.StdEncoding.EncodeToString(credentials),
		})
		if err != nil {
			return nil, err
		}
	}

	return tokens, nil
}
