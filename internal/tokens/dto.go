package tokens

import (
	"encoding/json"
	"time"
)

type Token struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	EnvID     string    `json:"env_id,omitempty"`
	Expiry    time.Time `json:"expiry,omitempty"`
	Key       string    `json:"key,omitempty"`
	Hash      string    `json:"hash,omitempty"`
	Name      string    `json:"name,omitempty"`
}

// IsExpired checks whether the token is expired or not.
func (t *Token) IsExpired() bool {
	return t.Expiry.Before(time.Now())
}

type CreateOptions struct {
	OrgKey []byte
	EnvID  string
	Expiry time.Duration
	Name   string `json:"name,omitempty"`
}

type CreateGraphQLOptions struct {
	EnvID  string
	Expiry time.Time
	Name   string
	Key    []byte
	Hash   string
}

type DecryptResponse struct {
	OrgKey []byte
	EnvID  string
	Expiry time.Time
	Name   string
}

type GetGraphQLOptions struct {
	Hash string `json:"hash"`
}

type ListOptions struct {
	EnvID string `json:"env_id,omitempty"`
}

// Custom marshaller for list options/filters.
func (o *ListOptions) MarshalJSON() ([]byte, error) {

	data := make(map[string]interface{})
	if o.EnvID != "" {
		data["env_id"] = map[string]interface{}{
			"_eq": o.EnvID,
		}
	}

	return json.Marshal(data)
}
