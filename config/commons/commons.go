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

type ConfigType string

const (
	ProjectConfig ConfigType = "ProjectConfig"
	AccountConfig ConfigType = "AccountConfig"
)
