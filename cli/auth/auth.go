package auth

import (
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
)

// To logout the user, simply delete account config
func Logout() error {
	return config.GetService().Delete(configCommons.AccountConfig)
}

func IsLoggedIn() bool {
	return config.GetService().Exists(configCommons.AccountConfig)
}
