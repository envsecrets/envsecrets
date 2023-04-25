package commons

type PermissionLevel string

const (
	OrgnisationLevelPermission PermissionLevel = "OrgnisationLevelPermission"
	ProjectLevelPermission     PermissionLevel = "ProjectLevelPermission"
	EnvironmentLevelPermission PermissionLevel = "EnvironmentLevelPermission"
)
