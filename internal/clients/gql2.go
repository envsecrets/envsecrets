package clients

import (
	"context"
	"os"

	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GQLClient2 struct {
	*graphql.Client
	BaseURL       string
	Authorization string
	log           *logrus.Logger
}

func NewGQLClient2(config *GQLConfig) *GQLClient2 {

	var response GQLClient2

	if config == nil {
		return &response
	}

	response.BaseURL = config.BaseURL
	response.Authorization = config.Authorization

	switch config.Type {
	case HasuraClientType:
		response.BaseURL = os.Getenv(string(NHOST_GRAPHQL_URL))
	}

	config.CustomHeaders = append(config.CustomHeaders, CustomHeader{
		Key:   "content-type",
		Value: "application/json",
	})

	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: response.Authorization,
			TokenType:   "Bearer",
		},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient(response.BaseURL, httpClient)
	response.Client = client

	if config.Logger != nil {
		response.log = config.Logger
	} else {
		response.log = logrus.New()
	}

	return &response
}
