package clients

import (
	"os"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"

	"github.com/machinebox/graphql"
)

type GQLClient struct {
	*clients.GQLClient
	log *logrus.Logger
}

type GQLConfig struct {
	BaseURL       string
	Authorization string
	Logger        *logrus.Logger
}

func NewGQLClient(config *GQLConfig) *GQLClient {

	client := clients.NewGQLClient(&clients.GQLConfig{
		BaseURL:       os.Getenv(string(NHOST_GRAPHQL_URL)),
		Authorization: config.Authorization,
	})

	response := GQLClient{
		GQLClient: client,
	}

	if config.Logger != nil {
		response.log = config.Logger
	} else {
		response.log = logrus.New()
	}

	return &response
}

func (c *GQLClient) Do(ctx context.ServiceContext, req *graphql.Request, resp interface{}) error {

	//	Parse the error
	if err := c.GQLClient.Do(ctx, req, &resp); err != nil {
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

			return c.GQLClient.Do(ctx, req, &resp)

		default:
			return apiError.ToError()
		}
	}

	return nil
}
