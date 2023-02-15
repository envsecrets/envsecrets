package permissions

import (
	"github.com/envsecrets/envsecrets/internal/permissions/commons"
)

type Service interface {
	Add(interface{}, commons.PermissionLevel) error
	Update(commons.Permissions, commons.PermissionLevel) error
	Exists(commons.PermissionLevel) bool
	Delete(commons.PermissionLevel) error
}

type DefaultPermissionService struct{}

/* func (*DefaultPermissionService) Save(payload interface{}, permissionsType commons.PermissionLevel) error {
	switch permissionsType {
	case commons.ProjectPermission:

		permissions, ok := payload.(commons.Project)
		if !ok {
			return errors.New("failed type assertion to project permissions")
		}
		return project.Save(&permissions)

	case commons.AccountPermission:

		permissions, ok := payload.(commons.Account)
		if !ok {
			return errors.New("failed type assertion to account permissions")
		}
		return account.Save(&permissions)
	}
	return nil
}

func (*DefaultPermissionService) Load(permissionsType commons.PermissionLevel) (interface{}, error) {
	switch permissionsType {
	case commons.ProjectPermission:
		return project.Load()
	case commons.AccountPermission:
		return account.Load()
	}

	return nil, nil
}

func (*DefaultPermissionService) Delete(permissionsType commons.PermissionLevel) error {
	switch permissionsType {
	case commons.ProjectPermission:
		return project.Delete()
	case commons.AccountPermission:
		return account.Delete()
	}

	return nil
}

func (*DefaultPermissionService) Exists(permissionsType commons.PermissionLevel) bool {
	switch permissionsType {
	case commons.ProjectPermission:
		return project.Exists()
	case commons.AccountPermission:
		return account.Exists()
	}

	return false
}
*/
