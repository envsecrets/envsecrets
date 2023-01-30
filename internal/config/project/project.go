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
	PROJECT_DIR        = filepath.Dir(commons.EXECUTABLE)
	PROJECT_CONFIG_DIR = filepath.Join(PROJECT_DIR, ".envsecrets")
	PROJECT_CONFIG_LOC = filepath.Join(PROJECT_CONFIG_DIR, CONFIG_FILENAME)
)

//	Save the provided config in its default location in the project root.
func Save(config *commons.Project) error {

	//	Create the configuration directory, if it doesn't already exist
	if err := os.MkdirAll(PROJECT_CONFIG_DIR, os.ModePerm); err != nil {
		return err
	}

	//	Marshal the yaml
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	//	Save the config file
	return ioutil.WriteFile(PROJECT_CONFIG_LOC, data, 0644)
}

//	Load, parse and return the available project config.
func Fetch() (*commons.Project, error) {

	//	Read the file
	data, err := ioutil.ReadFile(PROJECT_CONFIG_LOC)
	if err != nil {
		return nil, err
	}

	var config commons.Project

	//	Unmarshal its contents
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
