package users

import "time"

type User struct {
	ID          string    `json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Name        string    `json:"display_name,omitempty"`
	Email       string    `json:"email,omitempty"`
}
