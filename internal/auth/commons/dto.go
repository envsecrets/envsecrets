package commons

import "encoding/json"

type SigninOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResponse struct {
	MFA     interface{}            `json:"mfa"`
	Session map[string]interface{} `json:"session"`
	Keys    map[string]string      `json:"keys"`
}

type NhostSigninResponse struct {
	MFA     interface{}            `json:"mfa"`
	Session map[string]interface{} `json:"session"`
}

type NhostSession struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UpdatePasswordOptions struct {
	NewPassword string `json:"newPassword,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
}

func (o *UpdatePasswordOptions) Marshal() ([]byte, error) {
	return json.Marshal(o)
}
