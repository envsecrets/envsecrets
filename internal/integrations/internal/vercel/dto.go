package vercel

import "github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"

type SetupOptions struct {
	ConfigurationID string
	Source          string
	Next            string
	State           string
	OrgID           string
	Code            string
}

type PrepareCredentialsOptions struct {
	Code string
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

type ListOptions struct {
	Credentials map[string]interface{}
}

type SyncOptions struct {
	Credentials   map[string]interface{} `json:"credentials"`
	EntityDetails map[string]interface{} `json:"entity_details"`
	Data          *keypayload.KPMap      `json:"data"`
}

type ListProjectsResponse struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	ID                string       `json:"id,omitempty"`
	Name              string       `json:"name,omitempty"`
	Username          string       `json:"username,omitempty"`
	AccountID         string       `json:"accountId,omitempty"`
	LatestDeployments []Deployment `json:"latestDeployments,omitempty"`
}

type Deployment struct {
	Creator Creator `json:"creator"`
}

type Creator struct {
	UID      string `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
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
