package clients

import (
	"os"

	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
)

type GQLClient2 struct {
	*graphql.Client
	BaseURL       string
	Authorization string
	log           *logrus.Logger
}

func NewGQLClient2(config *GQLConfig) *GQLClient2 {

	var response GQLClient2

	response.BaseURL = config.BaseURL
	response.Authorization = config.Authorization

	switch config.Type {
	case HasuraClientType:
		response.BaseURL = os.Getenv(string(NHOST_GRAPHQL_URL))
	}

	httpClient := NewHTTPClient(&HTTPConfig{
		Authorization: response.Authorization,
		CustomHeaders: config.CustomHeaders,
		Headers:       config.Headers,
	})

	client := graphql.NewClient(response.BaseURL, httpClient)
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
