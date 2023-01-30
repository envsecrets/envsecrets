package client

import (
	"context"
	"net/http"

	accountConfig "github.com/envsecrets/envsecrets/config/account"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

var (
	GRAPHQL_CLIENT *graphql.Client
)

func init() {

	var httpClient *http.Client

	//	Load account config
	config, err := accountConfig.Load()
	if err == nil {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.AccessToken, TokenType: "Bearer"},
		)

		httpClient = oauth2.NewClient(context.Background(), src)
	}

	GRAPHQL_CLIENT = graphql.NewClient(NHOST_GRAPHQL_URL, httpClient)
}
