package graphql

import (
	"github.com/envsecrets/envsecrets/internal/secrets/internal/payload"
	"github.com/machinebox/graphql"
)

type SetOptions struct {
	EnvID   string                      `json:"env_id"`
	Data    map[string]*payload.Payload `json:"data"`
	Version *int                        `json:"version,omitempty"`
}

type GetOptions struct {
	Key     string `json:"key"`
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

// Returns the appropriate GraphQL query to use based on available options.
func (o *GetOptions) NewRequest() *graphql.Request {

	var req *graphql.Request

	var query string
	if o.Key != "" {

		if o.Version != nil {
			query = QUERY_BY_VERSION_BY_KEY
		} else {
			query = QUERY_LATEST_BY_KEY
		}

	} else {

		if o.Version != nil {
			query = QUERY_BY_VERSION_LATEST
		} else {
			query = QUERY_LATEST_ALL
		}
	}

	req = graphql.NewRequest(query)
	req.Var("env_id", o.EnvID)

	if o.Key != "" {
		req.Var("key", o.Key)
	}

	if o.Version != nil {
		req.Var("version", o.Version)
	}

	return req
}

type DeleteOptions struct {
	EnvID   string `json:"env_id"`
	Version int    `json:"version"`
}
