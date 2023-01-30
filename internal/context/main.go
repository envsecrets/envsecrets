package context

import (
	"context"

	accountConfig "github.com/envsecrets/envsecrets/config/account"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
)

type ServiceContext struct {
	context.Context
	Config *configCommons.Account
}

var DContext ServiceContext

func init() {

	//	Load account config
	config, _ := accountConfig.Load()
	DContext = ServiceContext{
		Context: context.Background(),
		Config:  config,
	}
}
