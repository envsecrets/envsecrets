package commons

import "github.com/envsecrets/envsecrets/internal/users"

type Project struct {
	Version     int    `json:"version,omitempty" yaml:"version,omitempty"`
	Workspace   string `json:"workspace,omitempty" yaml:"workspace,omitempty"`
	Project     string `json:"project,omitempty" yaml:"project,omitempty"`
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`
	Branch      string `json:"branch,omitempty" yaml:"branch,omitempty"`
}

type Account struct {
	AccessToken  string     `json:"access_token,omitempty" yaml:"accessToken,omitempty"`
	RefreshToken string     `json:"refresh_token,omitempty" yaml:"refreshToken,omitempty"`
	User         users.User `json:"user,omitempty" yaml:"user,omitempty"`
}
