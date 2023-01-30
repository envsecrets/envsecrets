package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/envsecrets/envsecrets/internal/config/commons"
	"gopkg.in/yaml.v3"
)

const (
	CONFIG_FILENAME = "config.yaml"
)

var (
	ACCOUNT_DIR        = filepath.Dir(commons.EXECUTABLE)
	ACCOUNT_CONFIG_DIR = filepath.Join(commons.HOME_DIR, ".envsecrets")
	ACCOUNT_CONFIG_LOC = filepath.Join(ACCOUNT_CONFIG_DIR, CONFIG_FILENAME)
)

//	Save the provided config in its default location in the root.
func Save(config *commons.Account) error {

	//	Create the configuration directory, if it doesn't already exist
	if err := os.MkdirAll(ACCOUNT_CONFIG_DIR, os.ModePerm); err != nil {
		return err
	}

	//	Marshal the yaml
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	//	Save the config file
	return ioutil.WriteFile(ACCOUNT_CONFIG_LOC, data, 0644)
}

//	Load, parse and return the available account config.
func Fetch() (*commons.Account, error) {

	//	Read the file
	data, err := ioutil.ReadFile(ACCOUNT_CONFIG_LOC)
	if err != nil {
		return nil, err
	}

	var config commons.Account

	//	Unmarshal its contents
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
