package contingency

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/envsecrets/envsecrets/cli/config/commons"
)

const (
	CONFIG_FILENAME = "contingency.json"
)

var (
	DIR        = filepath.Dir(commons.EXECUTABLE)
	CONFIG_DIR = filepath.Join(DIR, commons.CONFIG_FOLDER_NAME)
	CONFIG_LOC = filepath.Join(CONFIG_DIR, CONFIG_FILENAME)
)

//	Save the provided config in its default location in the root.
func Save(config *commons.Contingency) error {

	//	Load the existing Contingency data
	existing := make(commons.Contingency)

	secrets, err := Load()
	if err == nil {
		existing = *secrets
	}

	for key, payload := range *config {
		existing[key] = payload
	}

	//	Create the configuration directory, if it doesn't already exist
	if err := os.MkdirAll(CONFIG_DIR, os.ModePerm); err != nil {
		return err
	}

	//	Marshal the yaml
	data, err := json.MarshalIndent(existing, "", "\t")
	if err != nil {
		return err
	}

	//	Save the config file
	return ioutil.WriteFile(CONFIG_LOC, data, 0644)
}

//	Load, parse and return the available account config.
func Load() (*commons.Contingency, error) {

	//	Read the file
	data, err := ioutil.ReadFile(CONFIG_LOC)
	if err != nil {
		return nil, err
	}

	var config commons.Contingency

	//	Unmarshal its contents
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

//	Validate whether account config exists in file system or not
func Exists() bool {
	_, err := os.Stat(CONFIG_LOC)
	return !errors.Is(err, os.ErrNotExist)
}

func Delete() error {
	return os.Remove(CONFIG_LOC)
}