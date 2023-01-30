package config

import "os"

var (
	EXECUTABLE, _ = os.Executable()
	HOME_DIR, _   = os.UserHomeDir()

	//	Nhost Variables
	NHOST_AUTH_URL = os.Getenv("NHOST_AUTH_URL")
)
