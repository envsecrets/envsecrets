package commons

import (
	"encoding/json"
	"time"

	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/projects"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
)

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

func (p *Permissions) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

type CRUD struct {
	Create bool `json:"create,omitempty"`
	Read   bool `json:"read,omitempty"`
	Update bool `json:"update,omitempty"`
	Delete bool `json:"delete,omitempty"`
}

func (p *CRUD) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

type OrgnisationPermissions struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	OrgID        string                     `json:"org_id,omitempty"`
	Organisation organisations.Organisation `json:"organisation,omitempty"`

	User   userCommons.User `json:"user,omitempty"`
	UserID string           `json:"user_id,omitempty"`

	Key         string `json:"key,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}

//	Org's permissions structure will also have to be manually unmarshalled.
//	Because Hasura sends stringified JSON.
func (o *OrgnisationPermissions) GetPermissions() (*Permissions, error) {
	var response Permissions
	err := json.Unmarshal([]byte(o.Permissions), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type OrganisationPermissionsInsertOptions struct {
	OrgID  string `json:"org_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	RoleID string `json:"role_id,omitempty"`
	Key    string `json:"key,omitempty"`
}

type ProjectPermissions struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	ProjectID string           `json:"project_id,omitempty"`
	Project   projects.Project `json:"project,omitempty"`

	User   userCommons.User `json:"user,omitempty"`
	UserID string           `json:"user_id,omitempty"`

	Permissions string `json:"permissions,omitempty"`
}

func (p *ProjectPermissions) GetPermissions() (*Permissions, error) {
	var response Permissions
	err := json.Unmarshal([]byte(p.Permissions), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type ProjectPermissionsInsertOptions struct {
	ProjectID   string      `json:"project_id,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	Permissions Permissions `json:"permissions,omitempty"`
}

type EnvironmentPermissions struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	EnvID       string                   `json:"env_id,omitempty"`
	Environment environments.Environment `json:"environment,omitempty"`

	User   userCommons.User `json:"user,omitempty"`
	UserID string           `json:"user_id,omitempty"`

	Permissions string `json:"permissions,omitempty"`
}

type EnvironmentPermissionsInsertOptions struct {
	EnvID       string      `json:"env_id,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	Permissions Permissions `json:"permissions,omitempty"`
}

func (e *EnvironmentPermissions) GetPermissions() (*Permissions, error) {
	var response Permissions
	err := json.Unmarshal([]byte(e.Permissions), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
