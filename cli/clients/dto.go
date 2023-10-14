package clients

import "github.com/envsecrets/envsecrets/internal/users"

type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type LoginResponse struct {
	MFA struct {
		Ticket string `json:"ticket"`
	} `json:"mfa"`

	Session NhostSession `json:"session"`
}

type NhostSession struct {
	AccessToken          string     `json:"accessToken"`
	AccessTokenExpiresIn int        `json:"accessTokenExpiresIn"`
	RefreshToken         string     `json:"refreshToken"`
	User                 users.User `json:"user"`
}
