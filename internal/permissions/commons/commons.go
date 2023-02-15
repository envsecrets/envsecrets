package commons

import (
	"os"
)

var (
	EXECUTABLE, _ = os.Executable()
	HOME_DIR, _   = os.UserHomeDir()

	//	Nhost Variables
	NHOST_AUTH_URL = os.Getenv("NHOST_AUTH_URL")
)

type PermissionLevel string

const (
	ProjectLevelPermission     PermissionLevel = "ProjectLevelPermission"
	OrgnisationLevelPermission PermissionLevel = "OrgnisationLevelPermission"
)
