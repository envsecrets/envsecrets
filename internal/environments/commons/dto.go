package commons

import (
	"encoding/json"
	"time"

	intergrationCommons "github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
)

type Environment struct {
	ID        string    `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty" graphql:"name,omitempty"`
	ProjectID string    `json:"project_id,omitempty" graphql:"project_id"`
	UserID    string    `json:"user_id"`
}

func (w *Environment) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Environment) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
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
	IntegrationType intergrationCommons.IntegrationType `json:"integration_type,omitempty"`
	Password        string                              `json:"password" validate:"required"`
	Version         *int                                `json:"version,omitempty"`

	// Name of the secret to sync.
	Key string `json:"key,omitempty"`
}

type SyncRequestOptions struct {
	IntegrationType intergrationCommons.IntegrationType `json:"integration_type,omitempty"`
	Data            *keypayload.KPMap                   `json:"data"`
}

type SyncOptions struct {
	EnvID           string
	IntegrationType intergrationCommons.IntegrationType
	Secrets         *keypayload.KPMap
}
