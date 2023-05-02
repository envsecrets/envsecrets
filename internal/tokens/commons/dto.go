package commons

import "time"

type Token struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	EnvID     string    `json:"env_id,omitempty"`
	Expiry    time.Time `json:"expiry,omitempty"`
	Hash      string    `json:"hash,omitempty"`
	Name      string    `json:"name,omitempty"`
}

type CreateRequestOptions struct {
	EnvID  string `json:"env_id"`
	Expiry string `json:"expiry"`
	Name   string `json:"name,omitempty"`
}

type CreateServiceOptions struct {
	EnvID  string
	Expiry time.Duration
	Name   string `json:"name,omitempty"`
}

type CreateOptions struct {
	ID            string
	Key           []byte
	EnvID         string
	IssuedAt      time.Time
	NotBeforeTime time.Time
	Expiry        time.Time
	Name          string `json:"name,omitempty"`
}

type CreateGraphQLOptions struct {
	EnvID  string    `json:"env_id"`
	Expiry time.Time `json:"expiry"`
	Name   string    `json:"name,omitempty"`
	Hash   string    `json:"hash,omitempty"`
}

type GetGraphQLOptions struct {
	Hash string `json:"hash"`
}
