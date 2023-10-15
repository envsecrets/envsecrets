package heroku

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
)

// Prepares credentials to be saved in the database.
func PrepareCredentials(ctx context.ServiceContext, options *PrepareCredentialsOptions) (map[string]interface{}, error) {

	//	Exchange the code for Access Token
	response, err := GetAccessToken(ctx, &TokenRequestOptions{
		Code:        options.Code,
		RedirectURI: os.Getenv("REDIRECT_DOMAIN") + "/v1/integrations/heroku/callback/setup",
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token_type":    response.TokenType,
		"refresh_token": response.RefreshToken,
	}, nil
}

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Refresh access token
	access, err := RefreshToken(ctx, &TokenRefreshOptions{
		RefreshToken:  options.Credentials["refresh_token"].(string),
		OrgID:         options.OrgID,
		IntegrationID: options.IntegrationID,
	})
	if err != nil {
		return nil, err
	}

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: fmt.Sprintf("%s %s", access.TokenType, access.AccessToken),
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Accept",
				Value: "application/vnd.heroku+json; version=3",
			},
		},
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.heroku.com/apps", nil)
	if err != nil {
		return nil, err
	}

	var response ListProjectsResponse
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Refresh access token
	access, err := RefreshToken(ctx, &TokenRefreshOptions{
		RefreshToken:  options.Credentials["refresh_token"].(string),
		OrgID:         options.OrgID,
		IntegrationID: options.IntegrationID,
	})
	if err != nil {
		return err
	}

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Authorization: fmt.Sprintf("%s %s", access.TokenType, access.AccessToken),
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Accept",
				Value: "application/vnd.heroku+json; version=3",
			},
		},
	})

	//	Prepare the secrets.
	body, err := json.Marshal(options.Data.ToKVMap())
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.heroku.com/apps/%v/config-vars", options.EntityDetails["id"])
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if err := client.Run(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
