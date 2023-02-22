package context

import (
	"context"
)

type ServiceContext struct {
	context.Context
	Type ContextType
}

type Config struct {
	Type    ContextType
	Context context.Context
}

func NewContext(config *Config) ServiceContext {
	var response ServiceContext
	response.Context = config.Context
	if response.Context == nil {
		response.Context = context.Background()
	}
	response.Type = config.Type
	return response
}

//	var DContext ServiceContext

/* func init() {

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
*/
