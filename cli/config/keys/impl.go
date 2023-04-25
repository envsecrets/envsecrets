package keys

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/envsecrets/envsecrets/cli/config/commons"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG_FILENAME = "keys.yaml"
)

var (
	DIR        = filepath.Dir(commons.EXECUTABLE)
	CONFIG_DIR = filepath.Join(commons.HOME_DIR, commons.CONFIG_FOLDER_NAME)
	CONFIG_LOC = filepath.Join(CONFIG_DIR, CONFIG_FILENAME)
)

//	Save the provided config in its default location in the root.
func Save(config *commons.Keys) error {

	//	Create the configuration directory, if it doesn't already exist
	if err := os.MkdirAll(CONFIG_DIR, os.ModePerm); err != nil {
		return err
	}

	//	Marshal the yaml
	data, err := yaml.Marshal(config.Stringify())
	if err != nil {
		return err
	}

	//	Save the config file
	return ioutil.WriteFile(CONFIG_LOC, data, 0644)
}

//	Load, parse and return the available account config.
func Load() (*commons.Keys, error) {

	//	Read the file
	data, err := ioutil.ReadFile(CONFIG_LOC)
	if err != nil {
		return nil, err
	}

	var config commons.KeysStringified

	//	Unmarshal its contents
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Unstringify()
}

//	Validate whether account config exists in file system or not
func Exists() bool {
	_, err := os.Stat(CONFIG_LOC)
	return !errors.Is(err, os.ErrNotExist)
}

func Delete() error {
	return os.Remove(CONFIG_LOC)
}
