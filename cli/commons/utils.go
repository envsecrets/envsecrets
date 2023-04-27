package commons

import (
	"github.com/envsecrets/envsecrets/cli/config"
	"github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/labstack/gommon/log"
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
var ContingencyConfig *commons.Contingency

func Initialize() {

	//	Fetch the account config
	accountConfig, err := config.GetService().Load(commons.AccountConfig)
	if err != nil {
		log.Debug(err)
	} else {
		AccountConfig = accountConfig.(*commons.Account)
	}

	//	Fetch the project config
	projectConfig, err := config.GetService().Load(commons.ProjectConfig)
	if err != nil {
		log.Debug(err)
	} else {
		ProjectConfig = projectConfig.(*commons.Project)
	}

	//	Fetch the keys config
	keysConfig, err := config.GetService().Load(commons.KeysConfig)
	if err != nil {
		log.Debug(err)
	} else {
		KeysConfig = keysConfig.(*commons.Keys)
	}

	//	Fetch the Contingency config
	ContingencyConfig, err := config.GetService().Load(commons.ContingencyConfig)
	if err != nil {
		log.Debug(err)
	} else {
		ContingencyConfig = ContingencyConfig.(*commons.Contingency)
	}

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
