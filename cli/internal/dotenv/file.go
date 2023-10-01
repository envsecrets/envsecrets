package dotenv

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/dto"
)

const (
	FILENAME = ".env.local"
)

var (
	FILE_DIR = filepath.Join(commons.WORKING_DIR, commons.CONFIG_FOLDER_NAME)
	FILE_LOC = filepath.Join(FILE_DIR, FILENAME)
)

// Save the provided config in its default location in the root.
func Save(config *dto.KPMap) error {

	//	Load the existing Contingency data
	var existing dto.KPMap

	secrets, err := Load()
	if err == nil {
		existing = *secrets
	}

	existing.Overwrite(config)

	return Write(&existing)
}

// Writes the provided config in its default location in the root.
func Write(config *dto.KPMap) error {

	//	Create the configuration directory, if it doesn't already exist
	if err := os.MkdirAll(FILE_DIR, os.ModePerm); err != nil {
		return err
	}

	//	Marshal the content
	data, err := json.MarshalIndent(&config, "", "\t")
	if err != nil {
		return err
	}

	//	Save the config file
	return os.WriteFile(FILE_LOC, data, 0644)
}

// Load, parse and return the available account config.
func Load() (*dto.KPMap, error) {

	//	Read the file
	data, err := os.ReadFile(FILE_LOC)
	if err != nil {
		return nil, err
	}

	var config dto.KPMap

	//	Unmarshal its contents
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Calls the Load() function.
// But if that function call fails, it creates a new file and returns an empty config.
func LoadOrNew() (*dto.KPMap, error) {

	//	Read the file
	data, err := Load()
	if err != nil {

		//	If the local secret file doesn't exist, create a new one.
		if os.IsNotExist(err) {

			if err := Save(&dto.KPMap{}); err != nil {
				return nil, err
			}

			//	Re-call the function.
			return LoadOrNew()
		}

		return nil, err
	}

	return data, nil
}

// Validate whether account config exists in file system or not
func Exists() bool {
	_, err := os.Stat(FILE_LOC)
	return !errors.Is(err, os.ErrNotExist)
}

func Delete() error {
	return os.Remove(FILE_LOC)
}
