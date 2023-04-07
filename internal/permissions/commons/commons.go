package commons

type PermissionLevel string

const (
	RoleLevelPermission        PermissionLevel = "RoleLevelPermission"
	OrgnisationLevelPermission PermissionLevel = "OrgnisationLevelPermission"
	ProjectLevelPermission     PermissionLevel = "ProjectLevelPermission"
	EnvironmentLevelPermission PermissionLevel = "EnvironmentLevelPermission"
)
