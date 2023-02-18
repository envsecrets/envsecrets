package events

import (
	"net/http"

	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/permissions"
	permissionCommons "github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/labstack/echo/v4"
)

//	Called when a new row is inserted inside the `organisations` table.
func OrganisationInserted(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var organisation organisations.Organisation
	if err := MapToStruct(payload.Event.Data.New, &organisation); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	//	Insert all root permissions for the owner of the new organisation.
	service := permissions.GetService()

	if err := service.Insert(
		permissionCommons.OrgnisationLevelPermission,
		context.DContext,
		client.GRAPHQL_CLIENT,
		permissionCommons.OrganisationPermissionsInsertOptions{
			OrgID:  organisation.ID,
			UserID: payload.Event.SessionVariables.UserID,
			Permissions: permissionCommons.Permissions{
				PermissionsManage:  true,
				ProjectsManage:     true,
				EnvironmentsManage: true,
				SecretsWrite:       true,
			}}); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to insert root permissions for owner of this org",
			Error:   err.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "inserted root permissions for owner of this org",
	})
}

//	Called when a row is inserted/updated/deleted inside the `org_level_permissions` table.
func OrganisationLevelPermissions(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var organisation permissionCommons.OrgnisationPermissions
	if err := MapToStruct(payload.Event.Data.New, &organisation); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	incomingPermissions, err := organisation.GetPermissions()
	if err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal incoming permissions",
			Error:   err.Error(),
		})
	}

	//	Fetch the permissions service
	service := permissions.GetService()

	//	Initialize context with admin privileges
	ctx := context.DContext
	client := client.NewClient(&client.Config{AdminAccess: true})

	switch payload.Event.Op {
	case string(Insert):

		var permissions permissionCommons.Permissions

		//	If the user has been given permission to "manage projects" in the organisation,
		//	we have to give the user permission to manage every environment of every project.
		if incomingPermissions.ProjectsManage {
			permissions.EnvironmentsManage = true
		}

		//	If the user has been given permission to "write secrets" in the organisation,
		//	we have to give the user permission to write secrets in every environment of every project.
		if incomingPermissions.SecretsWrite {
			permissions.SecretsWrite = true
		}

		//	If the user has been given permission to "manage permissions" in the organisation,
		//	we have to give the user permission to manage permissions in every environment of every project.
		if incomingPermissions.PermissionsManage {
			permissions.PermissionsManage = true
		}

		//	Fetch all projects of the organisation
		projects, err := projects.List(ctx, client, &projects.ListOptions{
			OrgID: organisation.OrgID,
		})
		if err != nil {
			return c.JSON(http.StatusBadGateway, &APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to fetch projects for organisation",
				Error:   err.Message,
			})
		}

		//	Insert permissions for every project
		for _, item := range *projects {
			if err := service.Insert(
				permissionCommons.ProjectLevelPermission,
				ctx,
				client,
				permissionCommons.ProjectPermissionsInsertOptions{
					ProjectID:   item.ID,
					UserID:      organisation.UserID,
					Permissions: permissions}); err != nil {
				return c.JSON(http.StatusBadGateway, &APIResponse{
					Code:    http.StatusBadRequest,
					Message: "failed to insert permissions for project: " + item.ID,
					Error:   err.Message,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully inserted project level permissions",
	})
}

//	Called when a row is inserted/updated/deleted inside the `project_level_permissions` table.
func ProjectLevelPermissions(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	//	Unmarshal the data interface to our required entity.
	var project permissionCommons.ProjectPermissions
	if err := MapToStruct(payload.Event.Data.New, &project); err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal new data",
			Error:   err.Error(),
		})
	}

	incomingPermissions, err := project.GetPermissions()
	if err != nil {
		return c.JSON(http.StatusBadGateway, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to unmarshal incoming permissions",
			Error:   err.Error(),
		})
	}

	//	Fetch the permissions service
	service := permissions.GetService()

	//	Initialize context with admin privileges
	ctx := context.DContext
	client := client.NewClient(&client.Config{AdminAccess: true})

	switch payload.Event.Op {
	case string(Insert):

		var permissions permissionCommons.Permissions

		//	If the user has been given permission to "write secrets" in the project,
		//	we have to give the user permission to write secrets in every environment of every project.
		if incomingPermissions.SecretsWrite {
			permissions.SecretsWrite = true
		}

		//	If the user has been given permission to "manage permissions" in the project,
		//	we have to give the user permission to manage permissions in every environment of every project.
		if incomingPermissions.PermissionsManage {
			permissions.PermissionsManage = true
		}

		//	Fetch all projects of the organisation
		environments, err := environments.List(ctx, client, &environments.ListOptions{
			ProjectID: project.ProjectID,
		})
		if err != nil {
			return c.JSON(http.StatusBadGateway, &APIResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to fetch projects for project",
				Error:   err.Message,
			})
		}

		//	Insert permissions for every project
		for _, item := range *environments {
			if err := service.Insert(
				permissionCommons.EnvironmentLevelPermission,
				ctx,
				client,
				permissionCommons.EnvironmentPermissionsInsertOptions{
					EnvID:       item.ID,
					UserID:      project.UserID,
					Permissions: permissions}); err != nil {
				return c.JSON(http.StatusBadGateway, &APIResponse{
					Code:    http.StatusBadRequest,
					Message: "failed to insert permissions for environment: " + item.ID,
					Error:   err.Message,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, &APIResponse{
		Code:    http.StatusOK,
		Message: "successfully inserted environment level permissions",
	})
}

//	Called when a row is inserted/updated/deleted inside the `env_level_permissions` table.
func EnvironmentLevelPermissions(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload HasuraEventPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusOK, &APIResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse the body",
		})
	}

	return c.JSON(http.StatusBadRequest, &APIResponse{
		Code:    http.StatusBadRequest,
		Message: "un-built event endpoint",
	})
}
