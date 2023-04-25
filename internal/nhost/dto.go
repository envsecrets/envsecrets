package nhost

import "encoding/json"

type SignupOptions struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Options  interface{} `json:"options"`
}

func (o *SignupOptions) Marshal() ([]byte, error) {
	return json.Marshal(o)
}
