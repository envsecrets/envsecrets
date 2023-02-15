package commons

import (
	"time"

	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/projects"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
)

type Permissions struct {

	//	Invite/remove users and their permissions.
	PermissionsManage bool `json:"permissions_manage,omitempty"`

	//	Create/Delete projects.
	ProjectsManage bool `json:"projects_manage,omitempty"`

	//	Create/Delete environments.
	EnvironmentsManage bool `json:"environments_manage,omitempty"`

	//	Create/Update secrets.
	SecretsWrite bool `json:"secrets_write,omitempty"`
}

type OrgnisationPermissions struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	OrgID        string                     `json:"org_id,omitempty"`
	Organisation organisations.Organisation `json:"organisation,omitempty"`

	User   userCommons.User `json:"user,omitempty"`
	UserID string           `json:"user_id,omitempty"`

	Permissions Permissions `json:"permissions,omitempty"`
}

type ProjectPermissions struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	ProjectID string           `json:"project_id,omitempty"`
	Project   projects.Project `json:"project,omitempty"`

	User   userCommons.User `json:"user,omitempty"`
	UserID string           `json:"user_id,omitempty"`

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

	Permissions Permissions `json:"permissions,omitempty"`
}
