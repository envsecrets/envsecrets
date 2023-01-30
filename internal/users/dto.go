package users

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name      string    `json:"displayName,omitempty" graphql:"displayName,omitempty"`
	Email     string    `json:"email,omitempty" graphql:"email,omitempty"`
}

func (w *User) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *User) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}
