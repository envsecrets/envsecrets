package commons

import "encoding/json"

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
