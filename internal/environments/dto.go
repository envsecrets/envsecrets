package environments

import (
	"encoding/json"
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
	ProjectID string `json:"project_id,omitempty"`
}

// Custom marshaller for list options/filters.
func (o *ListOptions) MarshalJSON() ([]byte, error) {

	data := make(map[string]interface{})
	if o.ProjectID != "" {
		data["project_id"] = map[string]interface{}{
			"_eq": o.ProjectID,
		}
	}

	return json.Marshal(data)
}

type SyncOptions struct {
	EnvID    string
	EventIDs []string
	Pairs    *keypayload.KPMap
}
