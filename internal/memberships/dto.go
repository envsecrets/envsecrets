package memberships

import (
	"encoding/json"
	"time"

	usersCommons "github.com/envsecrets/envsecrets/internal/users/commons"
)

type Membership struct {
	ID        string            `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt time.Time         `json:"created_at,omitempty" graphql:"created_at"`
	UpdatedAt time.Time         `json:"updated_at,omitempty" graphql:"updated_at"`
	UserID    string            `json:"user_id,omitempty" graphql:"user_id"`
	User      usersCommons.User `json:"user,omitempty" graphql:"user"`
	OrgID     string            `json:"org_id,omitempty" graphql:"org_id"`
}

func (w *Membership) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Membership) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	OrgID string `json:"org_id,omitempty" graphql:"org_id"`
}

type ListOptions struct {
	OrgID string `json:"org_id,omitempty" graphql:"org_id"`
}

type CreateResponse struct {
	ID     string `json:"id,omitempty" graphql:"id,omitempty"`
	OrgID  string `json:"org_id,omitempty" graphql:"org_id,omitempty"`
	UserID string `json:"user_id,omitempty" graphql:"user_id,omitempty"`
}

type UpdateOptions struct {
}
