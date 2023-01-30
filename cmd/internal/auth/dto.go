package auth

import "github.com/envsecrets/envsecrets/internal/users"

type LoginResponse struct {
	MFA struct {
		Ticket string `json:"ticket"`
	} `json:"mfa"`

	Session struct {
		AccessToken          string     `json:"accessToken"`
		AccessTokenExpiresIn int        `json:"accessTokenExpiresIn"`
		RefreshToken         string     `json:"refreshToken"`
		User                 users.User `json:"user"`
	} `json:"session"`
}
