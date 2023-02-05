package invites

import (
	"encoding/json"
	"time"

	"github.com/envsecrets/envsecrets/internal/organisations"
	usersCommons "github.com/envsecrets/envsecrets/internal/users/commons"
)

type Invite struct {
	ID             string                     `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt      time.Time                  `json:"created_at,omitempty" graphql:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at,omitempty" graphql:"updated_at"`
	SenderID       string                     `json:"sender_id,omitempty" graphql:"sender_id"`
	UserBySenderID usersCommons.User          `json:"userBySenderId,omitempty" graphql:"userBySenderId"`
	ReceiverEmail  string                     `json:"receiver_email,omitempty" graphql:"receiver_email"`
	OrgID          string                     `json:"org_id,omitempty" graphql:"org_id"`
	Organisation   organisations.Organisation `json:"organisation,omitempty" graphql:"organisation"`
	Scope          string                     `json:"scope,omitempty" graphql:"scope"`
	Accepted       bool                       `json:"accepted,omitempty" graphql:"accepted"`
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
}
