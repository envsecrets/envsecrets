package commons

import (
	"github.com/envsecrets/envsecrets/config"
	"github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
)

//	Initialize common GQL Client for the CLI
var GQLClient = clients.NewGQLClient(&clients.GQLConfig{
	Type: clients.HasuraClientType,
})

//	Initialize common HTTP Client for the CLI
var HTTPClient *clients.HTTPClient

//	Initialize common context for the CLI
var DefaultContext = context.NewContext(&context.Config{Type: context.CLIContext})

func init() {

	//	Fetch the account config
	accountConfig, _ := config.GetService().Load(commons.AccountConfig)
	config := accountConfig.(*commons.Account)

	//	Initalize the HTTP client with bearer token from account config
	HTTPClient = clients.NewHTTPClient(&clients.HTTPConfig{
		BaseURL: API + "/v1",
	})

	if config != nil {
		HTTPClient.Authorization = "Bearer " + config.AccessToken
	}
}
