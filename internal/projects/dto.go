package projects

import (
	"encoding/json"
	"time"
)

type Project struct {
	ID        string    `json:"id" graphql:"id"`
	CreatedAt time.Time `json:"created_at,omitempty" graphql:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" graphql:"updated_at"`
	Name      string    `json:"name,omitempty" graphql:"name"`
	OrgID     string    `json:"org_id,omitempty" graphql:"org_id"`
}

func (w *Project) Marshal() ([]byte, error) {
	return json.Marshal(&w)
}

func (w *Project) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &w)
}

type CreateOptions struct {
	OrgID string `graphql:"org_id" json:"org_id"`
	Name  string `graphql:"name" json:"name"`
}

type CreateResponse struct {
	ID    string `graphql:"id" json:"id"`
	Name  string `graphql:"name" json:"name"`
	OrgID string `graphql:"org_id" json:"org_id"`
}

type UpdateOptions struct {
	Name string `json:"name" graphql:"name"`
}

type ListOptions struct {
	OrgID string `graphql:"org_id" json:"org_id"`
}
