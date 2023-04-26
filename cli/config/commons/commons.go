package commons

import (
	"os"
)

var (
	EXECUTABLE, _ = os.Executable()
	HOME_DIR, _   = os.UserHomeDir()

	//	Nhost Variables
	NHOST_AUTH_URL string
)

type ConfigType string

const (
	ProjectConfig ConfigType = "ProjectConfig"
	AccountConfig ConfigType = "AccountConfig"
	KeysConfig    ConfigType = "KeysConfig"

	CONFIG_FOLDER_NAME = ".envs"
)
