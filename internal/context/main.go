package context

import (
	"context"

	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
)

type ServiceContext struct {
	context.Context
	Config *configCommons.Account
}

var DContext ServiceContext

func init() {

	//	Load account config
	response, _ := config.GetService().Load(configCommons.AccountConfig)

	config, ok := response.(*configCommons.Account)
	if !ok {
		panic("failed type conversion for account config")
	}

	DContext = ServiceContext{
		Context: context.Background(),
		Config:  config,
	}
}
