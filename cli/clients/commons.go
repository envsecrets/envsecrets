package clients

// URLs
var NHOST_GRAPHQL_URL string
var NHOST_AUTH_URL string
var API string

type CustomHeader struct {
	Key   string
	Value string
}

type Header string

const (

	//	Standard HTTP headers
	AuthorizationHeader Header = "Authorization"
	ContentTypeHeader   Header = "Content-Type"
	TokenHeader         Header = "x-envsecrets-token"
	OrgIDHeader         Header = "x-envsecrets-org-id"
	HasuraWebhookSecret Header = "X-Hasura-Webhook-Secret"

	//	Hasura headers
	XHasuraAdminSecretHeader Header = "x-hasura-admin-secret"
)

type ClientType string

const (
	HTTPClientType    ClientType = "HTTPClient"
	HasuraClientType  ClientType = "HasuraClient"
	GraphQLClientType ClientType = "GraphQLClient"
)
