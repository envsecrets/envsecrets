package clients

import (
	"os"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"

	"github.com/machinebox/graphql"
)

type GQLClient struct {
	*graphql.Client
	BaseURL       string
	Authorization string
	Headers       []Header
	CustomHeaders []CustomHeader
	log           *logrus.Logger
}

type GQLConfig struct {
	Type          ClientType
	BaseURL       string
	Authorization string
	Headers       []Header
	CustomHeaders []CustomHeader
	Logger        *logrus.Logger
}

func NewGQLClient(config *GQLConfig) *GQLClient {

	var response GQLClient

	if config == nil {
		return &response
	}

	response.Headers = config.Headers
	response.CustomHeaders = config.CustomHeaders
	response.BaseURL = config.BaseURL
	response.Authorization = config.Authorization

	switch config.Type {
	case HasuraClientType:
		response.BaseURL = os.Getenv(string(NHOST_GRAPHQL_URL))
	}

	client := graphql.NewClient(response.BaseURL)
	response.Client = client

	if config.Logger != nil {
		response.log = config.Logger
	} else {
		response.log = logrus.New()
	}

	return &response
}

func (c *GQLClient) Do(ctx context.ServiceContext, req *graphql.Request, resp interface{}) error {

	//	Set Authorization Header
	if c.Authorization != "" {
		req.Header.Set(string(AuthorizationHeader), c.Authorization)
	}

	//	Set headers
	for _, item := range c.Headers {
		switch item {
		case XHasuraAdminSecretHeader:
			req.Header.Set(string(item), os.Getenv(string(NHOST_ADMIN_SECRET)))
		}
	}

	//	Set custom headers
	for _, item := range c.CustomHeaders {
		req.Header.Add(item.Key, item.Value)
	}

	//	Parse the error
	if err := c.Run(ctx, req, &resp); err != nil {
		apiError := ParseExternal(err)

		//	If it's a JWTExpired error,
		//	refresh the JWT and re-call the request.
		switch apiError.Type {
		case ErrorTypeJWTExpired:

			c.log.Debug("Request failed due to expired token. Refreshing access token to try again.")

			//	Fetch account configuration
			accountConfigPayload, err := config.GetService().Load(configCommons.AccountConfig)
			if err != nil {
				return New(err, "Failed to load account config", ErrorTypeInvalidAccountConfiguration, ErrorSourceSystem).ToError()
			}

			accountConfig := accountConfigPayload.(*configCommons.Account)

			authResponse, refreshErr := auth.RefreshToken(map[string]interface{}{
				"refreshToken": accountConfig.RefreshToken,
			})

			if refreshErr != nil {
				return New(refreshErr, "Failed to refresh access token", ErrorTypeInvalidToken, ErrorSourceNhost).ToError()
			}

			//	Save the refreshed account config
			refreshConfig := configCommons.Account{
				AccessToken:  authResponse.Session.AccessToken,
				RefreshToken: authResponse.Session.RefreshToken,
				User:         authResponse.Session.User,
			}

			if err := config.GetService().Save(refreshConfig, configCommons.AccountConfig); err != nil {
				return New(err, "Failed to save account config", ErrorTypeInvalidAccountConfiguration, ErrorSourceSystem).ToError()
			}

			//	Update the authorization header in client.
			if c.Authorization != "" {
				c.Authorization = "Bearer " + authResponse.Session.AccessToken
			}

			return c.Do(ctx, req, &resp)

		default:
			return apiError.ToError()
		}
	}

	return nil
}
