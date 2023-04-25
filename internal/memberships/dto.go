package memberships

import (
	"encoding/json"
	"time"
)

type Membership struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	UserID string `json:"user_id,omitempty"`
	OrgID  string `json:"org_id,omitempty"`
	RoleID string `json:"role_id,omitempty"`

	Key string `json:"key,omitempty"`
}

func (w *Membership) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Membership) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	UserID string `json:"user_id,omitempty"`
	OrgID  string `json:"org_id,omitempty"`
	RoleID string `json:"role_id,omitempty"`
	Key    string `json:"key,omitempty"`
}

type UpdateOptions struct {
	Name string `json:"name"`
}

type UpdateInviteLimitOptions struct {
	ID               string
	IncrementLimitBy int
}

type GetKeyOptions struct {
	UserID string `json:"user_id,omitempty"`
	OrgID  string `json:"org_id,omitempty"`
}
