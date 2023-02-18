package client

import "github.com/machinebox/graphql"

var (
	GRAPHQL_CLIENT *GQLClient
)

func init() {
	client := graphql.NewClient(NHOST_GRAPHQL_URL)
	GRAPHQL_CLIENT = &GQLClient{client, false}
}
