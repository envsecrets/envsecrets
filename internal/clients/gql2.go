package clients

import (
	"os"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GQLClient2 struct {
	*graphql.Client
	BaseURL       string
	Authorization *Authorization
	log           *logrus.Logger
}

type GQL2Config struct {
	Type          ClientType
	BaseURL       string
	Authorization *Authorization
	Logger        *logrus.Logger
}

func NewGQLClient2(config *GQL2Config) *GQLClient2 {

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

	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: response.Authorization.Token,
			TokenType:   string(response.Authorization.TokenType),
		},
	)

	httpClient := oauth2.NewClient(context.NewContext(&context.Config{Type: context.APIContext}), src)

	client := graphql.NewClient(response.BaseURL, httpClient)
	response.Client = client

	if config.Logger != nil {
		response.log = config.Logger
	} else {
		response.log = logrus.New()
	}

	return &response
}
