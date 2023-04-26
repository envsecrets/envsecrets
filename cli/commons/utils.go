package commons

import (
	"github.com/envsecrets/envsecrets/cli/config"
	"github.com/envsecrets/envsecrets/cli/config/commons"
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

//	Initialize configs
var AccountConfig *commons.Account
var ProjectConfig *commons.Project
var KeysConfig *commons.Keys

func Initialize() {

	//	Fetch the account config
	accountConfig, _ := config.GetService().Load(commons.AccountConfig)
	AccountConfig = accountConfig.(*commons.Account)

	//	Fetch the project config
	projectConfig, _ := config.GetService().Load(commons.ProjectConfig)
	ProjectConfig = projectConfig.(*commons.Project)

	//	Fetch the keys config
	keysConfig, _ := config.GetService().Load(commons.KeysConfig)
	KeysConfig = keysConfig.(*commons.Keys)

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

	if AccountConfig != nil {
		HTTPClient.Authorization = "Bearer " + AccountConfig.AccessToken
		GQLClient.Authorization = "Bearer " + AccountConfig.AccessToken
	}
}

func init() {
	Initialize()
}
