package organisations

import (
	"encoding/json"
	"time"
)

type Organisation struct {
	ID        string    `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty" graphql:"name,omitempty"`
	UserID    string    `json:"user_id,omitempty" graphql:"user_id,omitempty"`
	//	User      users.User `json:"user,omitempty" graphql:"user"`
}

func (w *Organisation) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Organisation) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	Name string `json:"name"`
}

type UpdateOptions struct {
	Name string `json:"name"`
}
