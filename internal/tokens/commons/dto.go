package commons

import "time"

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

type CreateRequestOptions struct {
	Password string `json:"password"`
	EnvID    string `json:"env_id"`
	Expiry   string `json:"expiry"`
	Name     string `json:"name,omitempty"`
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
