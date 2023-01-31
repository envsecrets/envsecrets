package memberships

import (
	"encoding/json"
	"time"
)

type Membership struct {
	ID          string    `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" graphql:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" graphql:"updated_at"`
	UserID      string    `json:"user_id,omitempty" graphql:"user_id"`
	WorkspaceID string    `json:"workspace_id,omitempty" graphql:"workspace_id"`
}

func (w *Membership) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Membership) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	WorkspaceID string `json:"workspace_id,omitempty" graphql:"workspace_id"`
}

type CreateResponse struct {
	ID          string `json:"id,omitempty" graphql:"id,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty" graphql:"workspace_id,omitempty"`
	UserID      string `json:"user_id,omitempty" graphql:"user_id,omitempty"`
}

type UpdateOptions struct {
}
