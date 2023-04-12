package commons

import "time"

type Token struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	EnvID     string    `json:"env_id,omitempty"`
	Expiry    time.Time `json:"expiry,omitempty"`
}

type CreateRequestOptions struct {
	EnvID  string `json:"env_id"`
	Expiry string `json:"expiry"`
}

type CreateServiceOptions struct {
	OrgID  string
	EnvID  string
	Expiry time.Duration
}

type CreateOptions struct {
	ID            string
	Key           []byte
	EnvID         string
	IssuedAt      time.Time
	NotBeforeTime time.Time
	Expiry        time.Time
}

type CreateGraphQLOptions struct {
	ID     string    `json:"id"`
	EnvID  string    `json:"env_id"`
	Expiry time.Time `json:"expiry"`
	Hash   string    `json:"hash"`
}

type GetGraphQLOptions struct {
	Hash string `json:"hash"`
}

type DecryptServiceOptions struct {
	OrgID string
	Token string
}

type DecryptOptions struct {
	Key   []byte
	Token string
}
