package commons

import (
	"time"

	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
)

type OauthAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type Integration struct {
	ID             string                 `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt      time.Time              `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	UserID         string                 `json:"user_id,omitempty" graphql:"user_id,omitempty"`
	OrgID          string                 `json:"org_id"`
	InstallationID string                 `json:"installation_id"`
	Type           IntegrationType        `json:"type"`
	Credentials    map[string]interface{} `json:"credentials,omitempty"`
}

type Integrations []Integration

type AddIntegrationOptions struct {

	//	Global
	OrgID          string          `json:"org_id"`
	InstallationID string          `json:"installation_id"`
	Type           IntegrationType `json:"type"`

	//	Especially for Vercel
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Scope       map[string]interface{} `json:"scope,omitempty"`
}

type ListIntegrationFilters struct {
	OrgID string          `json:"org_id"`
	Type  IntegrationType `json:"type"`
}

type Entity struct {
	ID             string          `json:"id"`
	Slug           string          `json:"slug"`
	ParentName     string          `json:"parent_name"`
	Name           string          `json:"name"`
	URL            string          `json:"url"`
	InstallationID string          `json:"installation_id"`
	Type           IntegrationType `json:"type"`
}

type Entities []Entity
type ListEntitiesRequest struct {
	OrgID string          `json:"org_id"`
	Type  IntegrationType `json:"type"`
}

type ListEntitiesRequestOptions struct {
	OrgID          string          `json:"org_id"`
	Type           IntegrationType `json:"type"`
	InstallationID string          `json:"installation_id"`
}

type SyncOptions struct {
	InstallationID string                           `json:"installation_id"`
	EntityDetails  map[string]interface{}           `json:"entity_details"`
	Credentials    map[string]interface{}           `json:"credentials"`
	Data           map[string]secretCommons.Payload `json:"data"`
}
