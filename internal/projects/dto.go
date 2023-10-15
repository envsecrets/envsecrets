package projects

import (
	"time"

	"github.com/envsecrets/envsecrets/internal/organisations"
)

type Project struct {
	ID           string                     `json:"id"`
	CreatedAt    time.Time                  `json:"created_at,omitempty"`
	UpdatedAt    time.Time                  `json:"updated_at,omitempty"`
	Name         string                     `json:"name,omitempty"`
	OrgID        string                     `json:"org_id,omitempty"`
	Organisation organisations.Organisation `json:"organisation,omitempty"`
	UserID       string                     `json:"user_id"`
}

type CreateOptions struct {
	OrgID string `json:"org_id"`
	Name  string `json:"name"`
}

type UpdateOptions struct {
	Name string `json:"name"`
}

type ListOptions struct {
	OrgID string `json:"org_id"`
}
