package projects

import (
	"encoding/json"
	"time"
)

type Project struct {
	ID          string    `json:"id" graphql:"id"`
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
	WorkspaceID string `graphql:"workspace_id" json:"workspace_id"`
	Name        string `graphql:"name" json:"name"`
}

type CreateResponse struct {
	ID          string `graphql:"id" json:"id"`
	Name        string `graphql:"name" json:"name"`
	WorkspaceID string `graphql:"workspace_id" json:"workspace_id"`
}

type UpdateOptions struct {
	Name string `json:"name" graphql:"name"`
}
