package commons

import (
	"encoding/base64"
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
	Salt         string `json:"salt,omitempty"`
}

func (k *Key) DecodePayload() (*Payload, error) {

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

	return &Payload{
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		ProtectedKey: protectedKey,
		Salt:         salt,
	}, nil
}

type Payload struct {
	PublicKey    []byte
	PrivateKey   []byte
	ProtectedKey []byte
	Salt         []byte
}

func (p *Payload) Validate() error {

	return nil
}

type CreateOptions struct {
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	ProtectedKey string `json:"protected_key"`
	Salt         string `json:"salt,omitempty"`
}

type CreateWithUserIDOptions struct {
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	ProtectedKey string `json:"protected_key"`
	Salt         string `json:"salt,omitempty"`
	UserID       string `json:"user_id,omitempty"`
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
}
