package events

import (
	"fmt"
	"time"

	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
)

type Type string

type Events []Event

type Event struct {
	ID            string                   `json:"id,omitempty"`
	CreatedAt     time.Time                `json:"created_at,omitempty"`
	UpdatedAt     time.Time                `json:"updated_at,omitempty"`
	Name          string                   `json:"name,omitempty"`
	IntegrationID string                   `json:"integration_id,omitempty"`
	EnvID         string                   `json:"env_id,omitempty"`
	EntityDetails map[string]interface{}   `json:"entity_details,omitempty"`
	Integration   integrations.Integration `json:"integration,omitempty"`
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
	case integrations.Railway:
		project := e.EntityDetails["project"].(map[string]interface{})
		return fmt.Sprintf("https://railway.app/project/%s/settings/variables", project["id"])
	case integrations.Hasura:
		return fmt.Sprintf("https://cloud.hasura.io/project/%s/env-vars", e.EntityDetails["id"])
	case integrations.Nhost:
		return fmt.Sprintf("https://app.nhost.io/%s/settings/secrets", e.EntityDetails["name"])
	case integrations.Heroku:
		return fmt.Sprintf("https://dashboard.heroku.com/apps/%s/settings", e.EntityDetails["name"])
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
	case integrations.Hasura:
		return e.EntityDetails["name"].(string)
	case integrations.Nhost:
		return e.EntityDetails["name"].(string)
	case integrations.Heroku:
		return e.EntityDetails["name"].(string)
	case integrations.Railway:
		project := e.EntityDetails["project"].(map[string]interface{})
		environment := e.EntityDetails["environment"].(map[string]interface{})
		return project["name"].(string) + "/" + environment["name"].(string)
	default:
		return ""
	}
}

// Get the type of entity by the type of it's integration.
func (e *Event) GetEntityType() string {
	switch e.Integration.Type {
	case integrations.Github:
		return "repository"
	case integrations.Vercel, integrations.CircleCI, integrations.Supabase, integrations.Railway, integrations.Hasura:
		return "project"
	case integrations.Gitlab:
		return "project/group"
	case integrations.ASM, integrations.GSM:
		return "secret"
	case integrations.Netlify:
		return "site"
	case integrations.Nhost, integrations.Heroku:
		return "app"
	default:
		return ""
	}
}

type ActionsGetOptions struct {
	EnvID string `json:"env_id,omitempty"`
}

type SyncOptions struct {
	ID string
	KP *keypayload.KPMap
}
