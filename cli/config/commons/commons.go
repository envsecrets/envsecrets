package commons

import (
	"os"
)

var (
	EXECUTABLE, _  = os.Executable()
	WORKING_DIR, _ = os.Getwd()
	HOME_DIR, _    = os.UserHomeDir()
)

type ConfigType string

const (
	ProjectConfig     ConfigType = "ProjectConfig"
	AccountConfig     ConfigType = "AccountConfig"
	KeysConfig        ConfigType = "KeysConfig"
	ContingencyConfig ConfigType = "ContingencyConfig"

	CONFIG_FOLDER_NAME = ".envs"
)
