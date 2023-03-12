package clients

import (
	"os"

	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"

	"github.com/machinebox/graphql"
)

type GQLClient struct {
	*graphql.Client
	BaseURL       string
	Authorization string
	Headers       []Header
	CustomHeaders []CustomHeader
}

type GQLConfig struct {
	Type          ClientType
	BaseURL       string
	Authorization string
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
	response.Authorization = config.Authorization

	switch config.Type {
	case HasuraClientType:
		response.BaseURL = os.Getenv(string(NHOST_GRAPHQL_URL))
	}

	client := graphql.NewClient(response.BaseURL)
	response.Client = client
	return &response
}

func (c *GQLClient) Do(ctx context.ServiceContext, req *graphql.Request, resp interface{}) *errors.Error {

	//	Set Authorization Header
	if c.Authorization != "" {
		req.Header.Add(string(AuthorizationHeader), c.Authorization)
	}

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

			if err := auth.RefreshAndSave(); err != nil {
				return errors.New(err, "failed to refresh and save auth token", errors.ErrorTypeBadResponse, errors.ErrorSourceGo)
			}

			//	Re-run the request
			c.Do(ctx, req, resp)

		} else {
			return apiError
		}
	}

	return nil
}
