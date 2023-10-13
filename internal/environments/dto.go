package environments

import (
	"time"

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

type SyncOptions struct {
	EnvID    string
	EventIDs []string
	Pairs    *keypayload.KPMap
}
