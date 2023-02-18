package client

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
	AdminAccess bool
}

type Config struct {
	AdminAccess bool
}

func NewClient(config *Config) *GQLClient {
	client := graphql.NewClient(NHOST_GRAPHQL_URL)
	return &GQLClient{client, config.AdminAccess}
}

func (c *GQLClient) Do(ctx context.ServiceContext, req *graphql.Request, resp interface{}) *errors.Error {

	//	If the context is being used by root,
	//	add the admin-secret header to request.
	//	Otherwise fetch the authorization header from saved account config.
	if c.AdminAccess {
		req.Header.Set(X_HASURA_ADMIN_SECRET, os.Getenv(NHOST_ADMIN_SECRET))
	} else {
		req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)
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

			response, refreshErr := auth.RefreshToken(map[string]interface{}{
				"refreshToken": ctx.Config.RefreshToken,
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

			//	Save the refresh credentials in current context
			ctx.Config = &refreshConfig

			//	Re-run the request
			c.Do(ctx, req, resp)
		} else {
			return apiError
		}
	}

	return nil
}
