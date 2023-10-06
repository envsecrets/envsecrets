package integrations

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

// Get the title of the integration by it's type.
func (i *Integration) GetTitle() string {
	switch i.Type {
	case Github:
		return "Github Actions"
	case Gitlab:
		return "Gitlab CI"
	case Vercel:
		return "Vercel"
	case ASM:
		return "AWS Secrets Manager"
	case GSM:
		return "Google Secrets Manager"
	case CircleCI:
		return "CircleCI"
	case Supabase:
		return "Supabase"
	case Netlify:
		return "Netlify"
	default:
		return ""
	}
}

// Get the subtitle of the integration by it's type.
func (i *Integration) GetSubtitle() string {
	switch i.Type {
	case Github:
		return "Your Github repository where we sync this environment's secrets."
	case Gitlab:
		return "Your Gitlab project/group where we sync this environment's secrets."
	case Vercel:
		return "Your Vercel project where we sync this environment's secrets."
	case ASM:
		return "Your ASM where we sync this environment's secrets."
	case GSM:
		return "Your GSM where we sync this environment's secrets."
	case CircleCI:
		return "Your CircleCI project where we sync this environment's secrets."
	case Supabase:
		return "Your Supabase project where we sync this environment's secrets."
	case Netlify:
		return "Your Netlify project where we sync this environment's secrets."
	default:
		return ""
	}
}

// Get the description of the integration by it's type.
func (i *Integration) GetDescription() string {
	switch i.Type {
	case Github:
		return "Make your secrets natively available in your repository's actions and workflows."
	case Gitlab:
		return "Make your secrets natively available in your repository's CI/CD pipelines."
	case Vercel:
		return "Make your secrets natively available in your project's environment variables."
	case ASM:
		return "Make your secrets natively available in your AWS Lambda functions."
	case GSM:
		return "Make your secrets natively available in your Google Cloud Functions."
	case CircleCI:
		return "Make your secrets natively available in your repository's CI/CD pipelines."
	case Supabase:
		return "Make your secrets natively available in your Supabase project's environment variables."
	case Netlify:
		return "Make your secrets natively available in your Netlify project's environment variables."
	default:
		return ""
	}
}

type Integrations []Integration

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
