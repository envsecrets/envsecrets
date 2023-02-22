package clients

import (
	"os"

	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"

	"github.com/machinebox/graphql"
)

type GQLClient struct {
	*graphql.Client
	BaseURL       string
	Headers       []Header
	CustomHeaders []CustomHeader
}

type GQLConfig struct {
	Type          ClientType
	BaseURL       string
	Headers       []Header
	CustomHeaders []CustomHeader
}

func NewGQLClient(config *GQLConfig) *GQLClient {

	var response GQLClient

	if config == nil {
		return &response
	}

	response.Headers = config.Headers
	response.CustomHeaders = config.CustomHeaders
	response.BaseURL = config.BaseURL

	switch config.Type {
	case HasuraClientType:
		response.BaseURL = os.Getenv(string(NHOST_GRAPHQL_URL))
	}

	client := graphql.NewClient(response.BaseURL)
	response.Client = client
	return &response
}

func (c *GQLClient) Do(ctx context.ServiceContext, req *graphql.Request, resp interface{}) *errors.Error {

	//	Set headers
	for _, item := range c.Headers {
		switch item {
		case XHasuraAdminSecretHeader:
			req.Header.Add(string(item), os.Getenv(string(NHOST_ADMIN_SECRET)))
		}
	}

	//	Set custom headers
	for _, item := range c.CustomHeaders {
		req.Header.Add(item.Key, item.Value)
	}

	return c.send(ctx, req, resp)
}

func (c *GQLClient) send(ctx context.ServiceContext, req *graphql.Request, resp interface{}) *errors.Error {

	//	Parse the error
	if err := c.Run(ctx, req, &resp); err != nil {

		apiError := errors.Parse(err)

		//	If it's a JWTExpired error,
		//	refresh the JWT and re-call the request.
		if apiError.IsType(errors.ErrorTypeJWTExpired) {

			//	Fetch account configuration
			accountConfigPayload, err := config.GetService().Load(configCommons.AccountConfig)
			if err != nil {
				return errors.New(err, "failed to load account configuration", errors.ErrorTypeInvalidAccountConfiguration, errors.ErrorSourceGo)
			}

			accountConfig := accountConfigPayload.(*configCommons.Account)

			response, refreshErr := auth.RefreshToken(map[string]interface{}{
				"refreshToken": accountConfig.RefreshToken,
			})

			if refreshErr != nil {
				return &errors.Error{
					Message: "failed to refresh auth token",
					Type:    errors.ErrorTypeTokenRefresh,
					Source:  errors.ErrorSourceNhost,
				}
			}

			//	Save the refreshed account config
			refreshConfig := configCommons.Account{
				AccessToken:  response.Session.AccessToken,
				RefreshToken: response.Session.RefreshToken,
				User:         response.Session.User,
			}

			if err := config.GetService().Save(refreshConfig, configCommons.AccountConfig); err != nil {
				return &errors.Error{
					Message: "failed to save refreshed login response",
					Type:    errors.ErrorTypeTokenRefresh,
					Source:  errors.ErrorSourceGo,
				}
			}

			//	Re-run the request
			c.Do(ctx, req, resp)
		} else {
			return apiError
		}
	}

	return nil
}
