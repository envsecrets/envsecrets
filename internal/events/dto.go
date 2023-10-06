package events

import (
	"fmt"
	"time"

	"github.com/envsecrets/envsecrets/internal/integrations"
)

type Type string

type Events []Event

type Event struct {
	ID            string                   `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt     time.Time                `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt     time.Time                `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name          string                   `json:"name,omitempty" graphql:"name,omitempty"`
	IntegrationID string                   `json:"integration_id,omitempty" graphql:"integration_id,omitempty"`
	EnvID         string                   `json:"env_id,omitempty" graphql:"env_id,omitempty"`
	EntityDetails map[string]interface{}   `json:"entity_details,omitempty" graphql:"entity_details,omitempty"`
	Integration   integrations.Integration `json:"integration,omitempty" graphql:"integration,omitempty"`
}

// Get the link of the entity link by the type of it's integration.
func (e *Event) GetEntityLink() string {
	switch e.Integration.Type {
	case integrations.Github:
		return fmt.Sprintf("https://github.com/%s/settings/secrets/actions", e.EntityDetails["full_name"])
	case integrations.Gitlab:
		return fmt.Sprintf("%s/-/settings/ci_cd", e.EntityDetails["web_url"])
	case integrations.Vercel:
		return fmt.Sprintf("https://vercel.com/%s/%s/settings/environment-variables", e.EntityDetails["username"], e.EntityDetails["name"])
	case integrations.ASM:
		return fmt.Sprintf("https://console.aws.amazon.com/secretsmanager/home?region=%s#/secret?name=%s", e.EntityDetails["region"], e.EntityDetails["name"])
	case integrations.GSM:
		return fmt.Sprintf("https://console.cloud.google.com/security/secret-manager/secret/%s/versions", e.EntityDetails["name"])
	case integrations.CircleCI:
		return fmt.Sprintf("https://app.circleci.com/settings/project/%s/environment-variables", e.EntityDetails["project_slug"])
		/*
			 	case integrations.Supabase:
					return fmt.Sprintf("https://app.supabase.io/project/%s/settings/secrets", e.EntityDetails["project_id"])
		*/
	case integrations.Netlify:
		return fmt.Sprintf("https://app.netlify.com/sites/%s/settings/env", e.EntityDetails["name"])
	default:
		return ""
	}
}

// Get the title of the entity by the type of it's integration.
func (e *Event) GetEntityTitle() string {
	switch e.Integration.Type {
	case integrations.Github:
		return e.EntityDetails["full_name"].(string)
	case integrations.Gitlab:
		return e.EntityDetails["name"].(string)
	case integrations.Vercel:
		return e.EntityDetails["username"].(string) + "/" + e.EntityDetails["name"].(string)
	case integrations.ASM:
		return e.EntityDetails["name"].(string)
	case integrations.GSM:
		return e.EntityDetails["name"].(string)
	case integrations.CircleCI:
		return e.EntityDetails["project_slug"].(string)
	case integrations.Supabase:
		return e.EntityDetails["name"].(string)
	case integrations.Netlify:
		return e.EntityDetails["name"].(string)
	default:
		return ""
	}
}

// Get the type of entity by the type of it's integration.
func (e *Event) GetEntityType() string {
	switch e.Integration.Type {
	case integrations.Github:
		return "repository"
	case integrations.Vercel, integrations.CircleCI, integrations.Supabase:
		return "project"
	case integrations.Gitlab:
		return "project/group"
	case integrations.ASM, integrations.GSM:
		return "secret"
	case integrations.Netlify:
		return "site"
	default:
		return ""
	}
}

type ActionsGetOptions struct {
	EnvID string `json:"env_id,omitempty" graphql:"env_id,omitempty"`
}
