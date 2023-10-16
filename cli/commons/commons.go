package commons

import (
	"github.com/envsecrets/envsecrets/cli/clients"
	"github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/dto"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/sirupsen/logrus"
)

var VERSION string

var Log = logrus.New()

// Initialize common GQL Client for the CLI
var GQLClient *clients.GQLClient

// Initialize common HTTP Client for the CLI
var HTTPClient *clients.HTTPClient

// Initialize common context for the CLI
var DefaultContext = context.NewContext(&context.Config{Type: context.CLIContext})

// Initialize configs
var AccountConfig *commons.Account
var ProjectConfig *commons.Project
var KeysConfig *commons.Keys

var Secret *dto.Secret
