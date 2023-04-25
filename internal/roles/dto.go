package roles

import (
	"encoding/json"
	"time"

	permissionCommons "github.com/envsecrets/envsecrets/internal/permissions/commons"
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

//	Permissions structure will also have to be manually unmarshalled.
//	Because Hasura sends stringified JSON.
func (o *Role) GetPermissions() (*permissionCommons.Permissions, error) {
	var response permissionCommons.Permissions
	err := json.Unmarshal([]byte(o.Permissions), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type RoleInsertOptions struct {
	OrgID       string                        `json:"org_id,omitempty"`
	Name        string                        `json:"name"`
	Permissions permissionCommons.Permissions `json:"permissions,omitempty"`
}
