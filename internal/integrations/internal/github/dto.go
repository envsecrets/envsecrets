package github

import (
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
)

type SetupOptions struct {
	InstallationID string
	SetupAction    string
	State          string
	OrgID          string
}

type SyncOptions struct {
	InstallationID string                 `json:"installation_id"`
	EntityDetails  map[string]interface{} `json:"entity_details"`
	Data           *keypayload.KPMap      `json:"data"`
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

type ListRepositoriesResponse struct {
	TotalCount           int          `json:"total_count"`
	RespositorySelection string       `json:"repository_selection"`
	Repositories         []Repository `json:"repositories"`
}

type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type ListOptions struct {
	InstallationID string
}
