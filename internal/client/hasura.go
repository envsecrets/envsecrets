package client

import (
	"github.com/machinebox/graphql"
)

var (
	GRAPHQL_CLIENT *graphql.Client
)

func init() {
	GRAPHQL_CLIENT = graphql.NewClient(NHOST_GRAPHQL_URL)
}
