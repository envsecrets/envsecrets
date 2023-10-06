package projects

import (
	"encoding/json"
	"time"
)

type Project struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty"`
	OrgID     string    `json:"org_id,omitempty"`
	UserID    string    `json:"user_id"`
}

func (w *Project) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Project) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	OrgID string `json:"org_id"`
	Name  string `json:"name"`
}

type UpdateOptions struct {
	Name string `json:"name"`
}

type ListOptions struct {
	OrgID string `json:"org_id"`
}
