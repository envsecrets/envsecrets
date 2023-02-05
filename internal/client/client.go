package client

import (
	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"

	"github.com/machinebox/graphql"
)

type GQLClient struct {
	*graphql.Client
}

func (c *GQLClient) Do(ctx context.ServiceContext, req *graphql.Request, resp interface{}) *errors.Error {

	//	Add Authorization token to the request
	req.Header.Set("Authorization", "Bearer "+ctx.Config.AccessToken)

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
