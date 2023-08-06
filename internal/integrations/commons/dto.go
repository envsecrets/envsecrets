package commons

import (
	"time"

	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
)

type OauthAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type Integration struct {
	ID             string          `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt      time.Time       `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	UserID         string          `json:"user_id,omitempty" graphql:"user_id,omitempty"`
	OrgID          string          `json:"org_id"`
	InstallationID string          `json:"installation_id"`
	Type           IntegrationType `json:"type"`
	Credentials    string          `json:"credentials,omitempty"`
}

type Integrations []Integration

type AddIntegrationOptions struct {

	//	Global
	OrgID          string          `json:"org_id"`
	InstallationID string          `json:"installation_id"`
	Type           IntegrationType `json:"type"`

	//	Especially for Vercel
	Credentials string                 `json:"credentials,omitempty"`
	Scope       map[string]interface{} `json:"scope,omitempty"`
}

type ListIntegrationFilters struct {
	OrgID string          `json:"org_id"`
	Type  IntegrationType `json:"type"`
}

type UpdateDetailsOptions struct {
	ID            string                 `json:"id"`
	EntityDetails map[string]interface{} `json:"entity_details"`
}

type UpdateCredentialsOptions struct {
	ID          string `json:"id"`
	Credentials string `json:"credentials"`
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

type SetupOptions struct {
	Options map[string]interface{} `json:"options"`
	OrgID   string                 `json:"org_id"`
}

type SyncOptions struct {
	EventID       string                 `json:"event_id"`
	IntegrationID string                 `json:"integration_id"`
	EntityDetails map[string]interface{} `json:"entity_details"`
	Data          *keypayload.KPMap      `json:"data"`
}
