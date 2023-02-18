package client

import "os"

var (
	NHOST_GRAPHQL_URL = os.Getenv("NHOST_GRAPHQL_URL")
)

const (

	//	HTTP headers
	AUTHORIZATION = "Authorization"

	//	Hasura headers
	X_HASURA_ADMIN_SECRET = "x-hasura-admin-secret"
)

const (

	//	Nhost constants
	NHOST_ADMIN_SECRET   = "NHOST_ADMIN_SECRET"
	NHOST_JWT_SECRET     = "NHOST_JWT_SECRET"
	NHOST_WEBHOOK_SECRET = "NHOST_WEBHOOK_SECRET"
	NHOST_AUTH_URL       = "NHOST_AUTH_URL"
	NHOST_STORAGE_URL    = "NHOST_STORAGE_URL"
	NHOST_FUNCTIONS_URL  = "NHOST_FUNCTIONS_URL"
)
