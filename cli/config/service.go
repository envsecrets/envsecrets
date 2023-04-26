package config

import (
	"errors"

	"github.com/envsecrets/envsecrets/cli/config/account"
	"github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/cli/config/keys"
	"github.com/envsecrets/envsecrets/cli/config/project"
)

type Service interface {
	Save(interface{}, commons.ConfigType) error
	Load(commons.ConfigType) (interface{}, error)
	Exists(commons.ConfigType) bool
	Delete(commons.ConfigType) error
}

type DefaultConfigService struct{}

func (*DefaultConfigService) Save(payload interface{}, configType commons.ConfigType) error {
	switch configType {
	case commons.ProjectConfig:

		config, ok := payload.(commons.Project)
		if !ok {
			return errors.New("failed type assertion to project config")
		}
		return project.Save(&config)

	case commons.AccountConfig:

		config, ok := payload.(commons.Account)
		if !ok {
			return errors.New("failed type assertion to account config")
		}
		return account.Save(&config)

	case commons.KeysConfig:

		config, ok := payload.(commons.Keys)
		if !ok {
			return errors.New("failed type assertion to keys config")
		}
		return keys.Save(&config)
	}
	return nil
}

func (*DefaultConfigService) Load(configType commons.ConfigType) (interface{}, error) {
	switch configType {
	case commons.ProjectConfig:
		return project.Load()
	case commons.AccountConfig:
		return account.Load()
	case commons.KeysConfig:
		return keys.Load()
	}

	return nil, nil
}

func (*DefaultConfigService) Delete(configType commons.ConfigType) error {
	switch configType {
	case commons.ProjectConfig:
		return project.Delete()
	case commons.AccountConfig:
		return account.Delete()
	case commons.KeysConfig:
		return keys.Delete()
	}

	return nil
}

func (*DefaultConfigService) Exists(configType commons.ConfigType) bool {
	switch configType {
	case commons.ProjectConfig:
		return project.Exists()
	case commons.AccountConfig:
		return account.Exists()
	case commons.KeysConfig:
		return keys.Exists()
	}

	return false
}
