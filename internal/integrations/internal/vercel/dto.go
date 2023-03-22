package vercel

type SetupOptions struct {
	ConfigurationID string
	Source          string
	Next            string
	State           string
	OrgID           string
	Token           string
	Code            string
}

type CodeExchangeResponse struct {
	TokenType      string `json:"token_type"`
	AccessToken    string `json:"access_token"`
	InstallationID string `json:"installation_id"`
	UserID         string `json:"user_id"`
	TeamID         string `json:"team_id"`
}

type InstallationAccessTokenResponse struct {
	Token                string `json:"token"`
	ExpiresAt            string `json:"expires_at"`
	RespositorySelection string `json:"repository_selection"`
	Permissions          struct {
		ActionsVariables  string `json:"actions_variables"`
		CodespacesSecrets string `json:"codespaces_secrets"`
		Deployments       string `json:"deployments"`
		Metadata          string `json:"metadata"`
		Secrets           string `json:"secrets"`
	} `json:"permissions"`
}

type RepositoryActionsSecretsPublicKeyResponse struct {
	Key   string `json:"key"`
	KeyID string `json:"key_id"`
}

type ListProjectsResponse struct {
	Projects []interface{} `json:"projects"`
}

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AccountID string `json:"accountId"`
}

type VercelResponse struct {
	Error map[string]interface{} `json:"error,omitempty"`
}

type VercelSecret struct {
	Error       map[string]interface{} `json:"error,omitempty"`
	ID          string                 `json:"uid"`
	Name        string                 `json:"name"`
	UserID      string                 `json:"userId"`
	TeamID      string                 `json:"teamId"`
	Decryptable bool                   `json:"decryptable"`
	Value       struct {
		Data interface{} `json:"data"`
		Type string      `json:"type"`
	} `json:"value"`
}
