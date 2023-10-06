package invites

import (
	"time"

	organisations "github.com/envsecrets/envsecrets/internal/organisations"
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

type CreateOptions struct {
	OrgID         string `json:"org_id,omitempty"`
	Scope         string `json:"scope,omitempty"`
	ReceiverEmail string `json:"receiver_email,omitempty"`
}

type ListOptions struct {
	//	OrgID          string                     `json:"org_id,omitempty"`
	Accepted bool `json:"accepted,omitempty"`
}

type CreateResponse struct {
	ID string `json:"id,omitempty"`
}

type UpdateOptions struct {
	Set SetUpdateOptions
}

type SetUpdateOptions struct {
	Accepted bool `json:"accepted,omitempty"`
}

type SendOptions struct {
	OrgID     string
	RoleID    string
	InviterID string

	//	Decrypted organisation key.
	Key          []byte
	InviteeEmail string
}

type InsertOptions struct {
	UserID string `json:"user_id,omitempty"`
	OrgID  string `json:"org_id,omitempty"`
	Key    string `json:"key,omitempty"`
	Email  string `json:"email,omitempty"`
	RoleID string `json:"role_id,omitempty"`
}
