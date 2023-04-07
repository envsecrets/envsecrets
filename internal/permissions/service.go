package permissions

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/envsecrets/envsecrets/internal/permissions/environment"
	"github.com/envsecrets/envsecrets/internal/permissions/organisation"
	"github.com/envsecrets/envsecrets/internal/permissions/project"
	"github.com/envsecrets/envsecrets/internal/permissions/roles"
)

type Service interface {
	Insert(commons.PermissionLevel, context.ServiceContext, *clients.GQLClient, interface{}) *errors.Error
	/*
		Update(commons.Permissions, commons.PermissionLevel) error
		Exists(commons.PermissionLevel) bool
		Count(commons.PermissionLevel) (int, error)
		Delete(commons.PermissionLevel) error
	*/
}

type DefaultPermissionService struct{}

func (*DefaultPermissionService) Insert(permissionsType commons.PermissionLevel, ctx context.ServiceContext, client *clients.GQLClient, options interface{}) *errors.Error {
	switch permissionsType {
	case commons.RoleLevelPermission:
		payload, ok := options.(commons.RoleInsertOptions)
		if !ok {
			return errors.New(nil, "failed type assertion to role level permissions", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
		}
		return roles.Insert(ctx, client, &payload)
	case commons.OrgnisationLevelPermission:
		payload, ok := options.(commons.OrganisationPermissionsInsertOptions)
		if !ok {
			return errors.New(nil, "failed type assertion to organisation level permissions", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
		}
		return organisation.Insert(ctx, client, &payload)
	case commons.ProjectLevelPermission:
		payload, ok := options.(commons.ProjectPermissionsInsertOptions)
		if !ok {
			return errors.New(nil, "failed type assertion to project level permissions", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
		}
		return project.Insert(ctx, client, &payload)
	case commons.EnvironmentLevelPermission:
		payload, ok := options.(commons.EnvironmentPermissionsInsertOptions)
		if !ok {
			return errors.New(nil, "failed type assertion to environment level permissions", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
		}
		return environment.Insert(ctx, client, &payload)
	}
	return nil
}

/*
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
