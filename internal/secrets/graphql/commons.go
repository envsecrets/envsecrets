package graphql

const (
	QUERY_LATEST_ALL = `query MyQuery($env_id: uuid!) {
		secrets(where: {env_id: {_eq: $env_id}}, order_by: {version: desc}, limit: 1) {
		  data
		  version
		}
	  }				  
	`

	QUERY_LATEST_BY_KEY = `query MyQuery($env_id: uuid!, $key: String!) {
		secrets(order_by: {version: desc}, limit: 1, where: {env_id: {_eq: $env_id}}) {
		  data(path: $key)
		  version
		}
	  }`

	QUERY_BY_VERSION_BY_KEY = `query MyQuery($env_id: uuid!, $key: String!, $version: Int!) {
		secrets(limit: 1, where: {env_id: {_eq: $env_id}, version: {_eq: $version}}) {
		  data(path: $key)
		  version
		}
	  }`

	QUERY_BY_VERSION_LATEST = `query MyQuery($env_id: uuid!, $version: Int!) {
		secrets(where: {env_id: {_eq: $env_id}, version: {_eq: $version}}) {
		  data
		  version
		}
	  }`
)
