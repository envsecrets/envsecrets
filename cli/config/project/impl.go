package project

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/envsecrets/envsecrets/cli/config/commons"
)

const (
	FILENAME = "config.json"
)

var (
	DIR        = commons.WORKING_DIR
	CONFIG_DIR = filepath.Join(DIR, commons.CONFIG_FOLDER_NAME)
	CONFIG_LOC = filepath.Join(CONFIG_DIR, FILENAME)
)

// Save the provided config in its default location in the project root.
func Save(config *commons.Project) error {

	//	Create the configuration directory, if it doesn't already exist
	if err := os.MkdirAll(CONFIG_DIR, os.ModePerm); err != nil {
		return err
	}

	//	Marshal the content
	data, err := json.MarshalIndent(&config, "", "\t")
	if err != nil {
		return err
	}

	//	Save the config file
	return os.WriteFile(CONFIG_LOC, data, 0644)
}

// Load, parse and return the available project config.
func Load() (*commons.Project, error) {

	//	Read the file
	data, err := os.ReadFile(CONFIG_LOC)
	if err != nil {
		return nil, err
	}

	var config commons.Project

	//	Unmarshal its contents
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate whether project config exists in file system or not
func Exists() bool {
	_, err := os.Stat(CONFIG_LOC)
	return !errors.Is(err, os.ErrNotExist)
}

func Delete() error {
	return os.Remove(CONFIG_LOC)
}
