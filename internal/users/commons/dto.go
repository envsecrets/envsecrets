package commons

import (
	"encoding/json"
	"time"
)

type User struct {
	ID          string    `json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Name        string    `json:"display_name,omitempty"`
	Email       string    `json:"email,omitempty"`
}

func (w *User) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *User) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}
