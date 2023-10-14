package clients

import (
	"os"

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

	if config == nil {
		return &response
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
		return err
	}

	return nil
}
