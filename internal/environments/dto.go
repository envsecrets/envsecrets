package environments

import (
	"time"

	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
)

type Environment struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty"`
	ProjectID string    `json:"project_id,omitempty"`
	UserID    string    `json:"user_id"`
}

type CreateOptions struct {
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
}

type UpdateOptions struct {
	Name string `json:"name"`
}

type ListOptions struct {
	ProjectID string `json:"project_id"`
}

type SyncWithPasswordRequestOptions struct {
	IntegrationType integrations.Type `json:"integration_type,omitempty"`
	Password        string            `json:"password" validate:"required"`
	Version         *int              `json:"version,omitempty"`

	// Name of the secret to sync.
	Key string `json:"key,omitempty"`
}

type SyncRequestOptions struct {
	IntegrationType integrations.Type `json:"integration_type,omitempty"`
	Data            *keypayload.KPMap `json:"data"`
}

type SyncOptions struct {
	EnvID           string
	IntegrationType integrations.Type
	Secrets         *keypayload.KPMap
}