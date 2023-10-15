package heroku

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
)

func GetAccessToken(ctx context.ServiceContext, options *TokenRequestOptions) (*TokenResponse, error) {

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "content-type",
				Value: "application/x-www-form-urlencoded",
			},
		},
	})

	//	Prepare the POST request data string.
	params := url.Values{}
	params.Add("client_secret", os.Getenv("HEROKU_CLIENT_SECRET"))

	if options.Code != "" {
		params.Add("grant_type", "authorization_code")
		params.Add("code", options.Code)
	} else if options.RefreshToken != "" {
		params.Add("grant_type", "refresh_token")
		params.Add("refresh_token", options.RefreshToken)
	} else {
		return nil, errors.New("either code or refresh token is required")
	}

	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://id.heroku.com/oauth/token", body)
	if err != nil {
		return nil, err
	}

	var response TokenResponse
	if err := client.Run(ctx, req, &response); err != nil {
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

		err = graphql.UpdateCredentials(ctx, client, &graphql.UpdateCredentialsOptions{
			ID:          options.IntegrationID,
			Credentials: base64.StdEncoding.EncodeToString(credentials),
		})
		if err != nil {
			return nil, err
		}
	}

	return tokens, nil
}
