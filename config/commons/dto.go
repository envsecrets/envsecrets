package commons

import userCommons "github.com/envsecrets/envsecrets/internal/users/commons"

type Account struct {
	AccessToken  string           `json:"access_token,omitempty" yaml:"accessToken,omitempty"`
	RefreshToken string           `json:"refresh_token,omitempty" yaml:"refreshToken,omitempty"`
	User         userCommons.User `json:"user,omitempty" yaml:"user,omitempty"`
}

type Project struct {
	Version      int    `json:"version,omitempty" yaml:"version,omitempty"`
	Organisation string `json:"organisation,omitempty" yaml:"organisation,omitempty"`
	Project      string `json:"project,omitempty" yaml:"project,omitempty"`
	Environment  string `json:"environment,omitempty" yaml:"environment,omitempty"`
	//	Branch       string `json:"branch,omitempty" yaml:"branch,omitempty"`
}
