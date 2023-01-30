package client

import (
	"github.com/hasura/go-graphql-client"
)

var (
	GRAPHQL_CLIENT *graphql.Client
)

func init() {
	GRAPHQL_CLIENT = graphql.NewClient(NHOST_GRAPHQL_URL, nil)
}
