package github

import (
	"fmt"

	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
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
	Secret         secretCommons.Secret   `json:"secret"`
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

func (r *Repository) ToEntity() *commons.Entity {
	return &commons.Entity{
		ID:         fmt.Sprint(r.ID),
		Slug:       r.FullName,
		URL:        r.HTMLURL,
		Type:       commons.Github,
		Name:       r.Name,
		ParentName: r.Owner.Login,
	}
}
