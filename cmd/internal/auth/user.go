package auth

import (
	accountConfig "github.com/envsecrets/envsecrets/config/account"
)

//	Validate whether user is logged in
func IsLoggedIn() bool {

	config, err := accountConfig.Load()
	if err != nil {
		return false
	}

	return len(config.AccessToken) > 0
}
