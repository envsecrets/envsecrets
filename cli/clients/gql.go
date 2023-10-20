package clients

import (
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
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
		BaseURL:       NHOST_GRAPHQL_URL,
		Authorization: config.Authorization,
		ErrorHandler:  getErrorHandler(config.Logger),
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
		return err
	}

	return nil
}

func getErrorHandler(log *logrus.Logger) func(*clients.GQLClient, error) error {

	return func(c *clients.GQLClient, err error) error {

		apiError := ParseExternal(err)

		//	If it's a JWTExpired error,
		//	refresh the JWT and re-call the request.
		switch apiError.Type {
		case ErrorTypeJWTExpired, ErrorTypeMalformedHeader:

			log.Debug("JWT Expired or Header Malformed. Refreshing...")

			//	Fetch account configuration
			accountConfigPayload, err := config.GetService().Load(configCommons.AccountConfig)
			if err != nil {
				return New(err, "Failed to load account config", ErrorTypeInvalidAccountConfiguration, ErrorSourceSystem).ToError()
			}

			accountConfig := accountConfigPayload.(*configCommons.Account)

			ctx := context.NewContext(&context.Config{Type: context.CLIContext})

			//	Initialize a new Nhost client.
			nhostClient := NewNhostClient(&NhostConfig{
				Logger: log,
			})

			authResponse, refreshErr := auth.GetService().RefreshToken(ctx, nhostClient.NhostClient, &auth.RefreshTokenOptions{
				RefreshToken: accountConfig.RefreshToken,
			})

			if refreshErr != nil {
				return New(refreshErr, "Failed to refresh access token", ErrorTypeInvalidToken, ErrorSourceNhost).ToError()
			}

			//	Save the refreshed account config
			refreshConfig := configCommons.Account{
				AccessToken:  authResponse.AccessToken,
				RefreshToken: authResponse.RefreshToken,
				User:         authResponse.User,
			}

			if err := config.GetService().Save(refreshConfig, configCommons.AccountConfig); err != nil {
				return err
			}

			//	Update the authorization header in client.
			c.Authorization = "Bearer " + authResponse.AccessToken

			//	Re-do the request
			c.RedoRequestOnError = true

		default:
			return apiError.ToError()
		}

		return nil
	}
}
