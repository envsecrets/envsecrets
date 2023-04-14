package commons

import (
	"encoding/json"
	"time"

	"github.com/envsecrets/envsecrets/internal/organisations"
)

type Invite struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	UserID string `json:"user_id,omitempty"`

	OrgID        string                     `json:"org_id,omitempty"`
	Organisation organisations.Organisation `json:"organisation,omitempty"`

	Key   string `json:"key,omitempty"`
	Email string `json:"email,omitempty"`

	RoleID string `json:"role_id,omitempty"`

	Accepted bool `json:"accepted,omitempty"`
}

func (w *Invite) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Invite) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	OrgID         string `json:"org_id,omitempty" graphql:"org_id,omitempty"`
	Scope         string `json:"scope,omitempty" graphql:"scope"`
	ReceiverEmail string `json:"receiver_email,omitempty" graphql:"receiver_email"`
}

type ListOptions struct {
	//	OrgID          string                     `json:"org_id,omitempty" graphql:"org_id"`
	Accepted bool `json:"accepted,omitempty" graphql:"accepted"`
}

type CreateResponse struct {
	ID string `json:"id,omitempty" graphql:"id,omitempty"`
}

type UpdateOptions struct {
	Accepted bool `json:"accepted,omitempty"`
}
