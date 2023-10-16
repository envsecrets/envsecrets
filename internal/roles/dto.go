package roles

import (
	"time"
)

type Role struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Name string `json:"name"`

	OrgID string `json:"org_id,omitempty"`
	//	Organisation organisations.Organisation `json:"organisation,omitempty"`

	Permissions string `json:"permissions,omitempty"`
}

type RoleInsertOptions struct {
	OrgID       string      `json:"org_id,omitempty"`
	Name        string      `json:"name"`
	Permissions Permissions `json:"permissions,omitempty"`
}

type Permissions struct {

	//	Other members and their permissions.
	Permissions CRUD `json:"permissions,omitempty"`

	//	Projects and their environments.
	Projects CRUD `json:"projects,omitempty"`

	//	Environments and their secrets.
	Environments CRUD `json:"environments,omitempty"`

	//	Add/Delete Integrations.
	Integrations CRUD `json:"integrations,omitempty"`
}

type CRUD struct {
	Create bool `json:"create,omitempty"`
	Read   bool `json:"read,omitempty"`
	Update bool `json:"update,omitempty"`
	Delete bool `json:"delete,omitempty"`
}
