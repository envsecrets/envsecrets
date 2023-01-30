package workspaces

import (
	"encoding/json"
	"time"
)

type Workspace struct {
	ID        string    `json:"id,omitempty" graphql:"id"`
	CreatedAt time.Time `json:"created_at,omitempty" graphql:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" graphql:"updated_at"`
	Name      string    `json:"name,omitempty" graphql:"name"`
	UserID    string    `json:"user_id,omitempty" graphql:"user_id"`
	//	User      users.User `json:"user,omitempty" graphql:"user"`
}

func (w *Workspace) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Workspace) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	Name string `json:"name,omitempty"`
}

type UpdateOptions struct {
	Name string `json:"name,omitempty"`
}
