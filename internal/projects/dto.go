package projects

import (
	"encoding/json"
	"time"
)

type Project struct {
	ID          string    `json:"id,omitempty" graphql:"id"`
	CreatedAt   time.Time `json:"created_at,omitempty" graphql:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" graphql:"updated_at"`
	Name        string    `json:"name,omitempty" graphql:"name"`
	WorkspaceID string    `json:"workspace_id,omitempty" graphql:"workspace_id"`
}

func (w *Project) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Project) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	Name string `json:"name,omitempty"`
}

type UpdateOptions struct {
	Name string `json:"name,omitempty"`
}
