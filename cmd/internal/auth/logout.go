package auth

import (
	"os"

	accountConfig "github.com/envsecrets/envsecrets/config/account"
)

//	To logout the user, simply delete account config
func Logout() error {
	return os.Remove(accountConfig.ACCOUNT_CONFIG_LOC)
}
