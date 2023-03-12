package context

import (
	"context"

	"github.com/labstack/echo/v4"
)

type ServiceContext struct {
	context.Context
	Type        ContextType
	EchoContext echo.Context
}

type Config struct {
	Type        ContextType
	Context     context.Context
	EchoContext echo.Context
}

func NewContext(config *Config) ServiceContext {
	var response ServiceContext
	response.Context = config.Context
	response.EchoContext = config.EchoContext
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
