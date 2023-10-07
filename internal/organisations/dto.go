package organisations

import (
	"time"
)

type Organisation struct {
	ID          string    `json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Name        string    `json:"name,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	InviteLimit *int      `json:"invite_limit,omitempty"`
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
