package commons

import (
	"encoding/base64"
	"time"

	"github.com/envsecrets/envsecrets/internal/errors"
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

func (k *Key) DecodePayload() (*Payload, *errors.Error) {

	salt, er := base64.StdEncoding.DecodeString(k.Salt)
	if er != nil {
		return nil, errors.New(er, "Failed to base64 decode salt", errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	protectedKey, er := base64.StdEncoding.DecodeString(k.ProtectedKey)
	if er != nil {
		return nil, errors.New(er, "Failed to base64 decode protection key", errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	privateKey, er := base64.StdEncoding.DecodeString(k.PrivateKey)
	if er != nil {
		return nil, errors.New(er, "Failed to base64 decode private key", errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	publicKey, er := base64.StdEncoding.DecodeString(k.PublicKey)
	if er != nil {
		return nil, errors.New(er, "Failed to base64 decode public key", errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
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

func (p *Payload) Validate() *errors.Error {

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
