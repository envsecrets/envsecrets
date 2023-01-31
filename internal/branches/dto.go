package branches

import (
	"encoding/json"
	"time"
)

type Branch struct {
	ID            string    `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name          string    `json:"name,omitempty" graphql:"name,omitempty"`
	EnvironmentID string    `json:"environment_id,omitempty" graphql:"environment_id"`
}

func (w *Branch) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Branch) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	Name          string `json:"name" graphql:"name"`
	EnvironmentID string `json:"environment_id" graphql:"environment_id"`
}

type CreateResponse struct {
	ID   string `json:"id,omitempty" graphql:"id,omitempty"`
	Name string `json:"name,omitempty" graphql:"name,omitempty"`
}

type UpdateOptions struct {
	Name string `json:"name"`
}
