package organisations

import (
	"encoding/json"
	"time"
)

type Organisation struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty"`
	UserID    string    `json:"user_id,omitempty"`

	ServerKey string `json:"server_copy,omitempty"`
}

func (w *Organisation) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Organisation) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	Name   string `json:"name"`
	UserID string `json:"user_id,omitempty"`
}

type UpdateOptions struct {
	Name string `json:"name"`
}

type UpdateInviteLimitOptions struct {
	ID               string
	IncrementLimitBy int
}

type UpdateServerKeyCopyOptions struct {
	OrgID string
	Key   string
}
