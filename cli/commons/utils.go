package commons

import (
	"github.com/envsecrets/envsecrets/config"
	"github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"
)

//	Initialize common GQL Client for the CLI
var GQLClient *clients.GQLClient

//	Initialize common HTTP Client for the CLI
var HTTPClient *clients.HTTPClient

//	Initialize common context for the CLI
var DefaultContext = context.NewContext(&context.Config{Type: context.CLIContext})

var Logger = logrus.New()

func init() {

	//	Fetch the account config
	accountConfig, _ := config.GetService().Load(commons.AccountConfig)
	config := accountConfig.(*commons.Account)

	//	Initalize the HTTP client with bearer token from account config
	HTTPClient = clients.NewHTTPClient(&clients.HTTPConfig{
		BaseURL: API + "/v1",
		Logger:  Logger,
	})

	//	Initialize GQL client
	GQLClient = clients.NewGQLClient(&clients.GQLConfig{
		BaseURL: NHOST_GRAPHQL_URL,
		Logger:  Logger,
	})

	if config != nil {
		HTTPClient.Authorization = "Bearer " + config.AccessToken
		GQLClient.Authorization = "Bearer " + config.AccessToken
	}
}
