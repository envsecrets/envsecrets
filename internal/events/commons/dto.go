package commons

import (
	"time"

	integrationsCommons "github.com/envsecrets/envsecrets/internal/integrations/commons"
)

type Type string

type Events []Event

type Event struct {
	ID            string                          `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt     time.Time                       `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt     time.Time                       `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name          string                          `json:"name,omitempty" graphql:"name,omitempty"`
	IntegrationID string                          `json:"integration_id,omitempty" graphql:"integration_id,omitempty"`
	EnvID         string                          `json:"env_id,omitempty" graphql:"env_id,omitempty"`
	EntityDetails map[string]interface{}          `json:"entity_details,omitempty" graphql:"entity_details,omitempty"`
	Integration   integrationsCommons.Integration `json:"integration,omitempty" graphql:"integration,omitempty"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
