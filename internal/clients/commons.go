package clients

type Variable string

const (

	//	API constants
	API Variable = "API"

	//	Nhost constants
	NHOST_ADMIN_SECRET   Variable = "NHOST_ADMIN_SECRET"
	NHOST_JWT_SECRET     Variable = "NHOST_JWT_SECRET"
	NHOST_WEBHOOK_SECRET Variable = "NHOST_WEBHOOK_SECRET"
	NHOST_AUTH_URL       Variable = "NHOST_AUTH_URL"
	NHOST_STORAGE_URL    Variable = "NHOST_STORAGE_URL"
	NHOST_FUNCTIONS_URL  Variable = "NHOST_FUNCTIONS_URL"
	NHOST_GRAPHQL_URL    Variable = "NHOST_GRAPHQL_URL"
)

type ClientType string

const (
	HTTPClientType    ClientType = "HTTPClient"
	GithubClientType  ClientType = "GithubClient"
	VaultClientType   ClientType = "VaultClient"
	HasuraClientType  ClientType = "HasuraClient"
	GraphQLClientType ClientType = "GraphQLClient"
)

type CustomHeader struct {
	Key   string
	Value string
}

type Header string

const (

	//	Standard HTTP headers
	AuthorizationHeader Header = "Authorization"
	ContentTypeHeader   Header = "Content-Type"

	//	Hasura headers
	XHasuraAdminSecretHeader Header = "x-hasura-admin-secret"

	//	Github headers
	AcceptHeader Header = "Accept"

	VaultTokenHeader Header = "X-Vault-Token"
)

type CallType string

const (
	HTTPCall    CallType = "HTTPCall"
	GraphQLCall CallType = "GraphQLCall"
)
