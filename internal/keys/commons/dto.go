package commons

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Key struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserID    string    `json:"user_id,omitempty"`

	PublicKey    string `json:"public_key,omitempty"`
	PrivateKey   string `json:"private_key,omitempty"`
	ProtectedKey string `json:"protected_key,omitempty"`
	SyncKey      string `json:"sync_key,omitempty"`
	Salt         string `json:"salt,omitempty"`
}

func (k *Key) Decode() (*Payload, error) {

	salt, err := base64.StdEncoding.DecodeString(k.Salt)
	if err != nil {
		return nil, err
	}

	protectedKey, err := base64.StdEncoding.DecodeString(k.ProtectedKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := base64.StdEncoding.DecodeString(k.PrivateKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := base64.StdEncoding.DecodeString(k.PublicKey)
	if err != nil {
		return nil, err
	}

	var syncKey []byte
	if k.SyncKey != "" {
		syncKey, err = base64.StdEncoding.DecodeString(k.SyncKey)
		if err != nil {
			return nil, err
		}
	}

	return &Payload{
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		ProtectedKey: protectedKey,
		SyncKey:      syncKey,
		Salt:         salt,
	}, nil
}

type Payload struct {
	PublicKey    []byte
	PrivateKey   []byte
	ProtectedKey []byte
	SyncKey      []byte
	Salt         []byte
}

type DecryptOptions struct {
	Password string
	OrgID    string
}

func (o *DecryptOptions) Marshal() ([]byte, error) {
	return json.Marshal(o)
}

type CreateOptions struct {
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	ProtectedKey string `json:"protected_key"`
	SyncKey      string `json:"sync_key,omitempty"`
	Salt         string `json:"salt,omitempty"`
}

type CreateWithUserIDOptions struct {
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	ProtectedKey string `json:"protected_key"`
	Salt         string `json:"salt,omitempty"`
	SyncKey      string `json:"sync_key,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}

type CreateSyncKeyOptions struct {
	SyncKey string `json:"sync_key,omitempty"`
}

type GetPublicKeyOptions struct {
	Email  string `query:"email,omitempty"`
	UserID string `query:"user_id,omitempty"`
}

type IssueKeyPairResponse struct {
	PublicKey           []byte `json:"public_key"`
	PrivateKey          []byte `json:"private_key"`
	DecryptedPrivateKey []byte `json:"decrypted_private_key"`
	ProtectedKey        []byte `json:"protected_key"`
	Salt                []byte `json:"salt,omitempty"`
	SyncKey             []byte `json:"sync_key,omitempty"`
}
