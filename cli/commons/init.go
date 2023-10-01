package commons

import (
	"github.com/envsecrets/envsecrets/cli/config"
	"github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/sirupsen/logrus"
)

func Initialize(log *logrus.Logger) {

	if log == nil {
		log = logrus.New()
	}

	//	Fetch the account config
	accountConfig, err := config.GetService().Load(commons.AccountConfig)
	if err != nil {
		log.Debug(err)
	} else {
		AccountConfig = accountConfig.(*commons.Account)
	}

	//	Initalize the HTTP client with bearer token from account config
	HTTPClient = clients.NewHTTPClient(&clients.HTTPConfig{
		Type:    clients.HasuraClientType,
		BaseURL: API + "/v1",
		Logger:  log,
	})

	//	Initialize GQL client
	GQLClient = clients.NewGQLClient(&clients.GQLConfig{
		BaseURL: NHOST_GRAPHQL_URL,
		Logger:  log,
	})

	if AccountConfig != nil {
		HTTPClient.Authorization = "Bearer " + AccountConfig.AccessToken
		GQLClient.Authorization = "Bearer " + AccountConfig.AccessToken
	}

	//	Fetch the keys config
	keysConfig, err := config.GetService().Load(commons.KeysConfig)
	if err != nil {
		log.Debug(err)
	} else {
		KeysConfig = keysConfig.(*commons.Keys)
	}
}

/* func init() {
	Initialize()
}
*/
