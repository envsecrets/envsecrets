package commons

import (
	"encoding/json"
	"fmt"
	"time"
)

type Type string

type Secrets []Secret

type Secret struct {
	ID        string             `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name      string             `json:"name,omitempty" graphql:"name,omitempty"`
	UserID    string             `json:"user_id,omitempty" graphql:"user_id,omitempty"`
	EnvID     string             `json:"env_id,omitempty" graphql:"env_id,omitempty"`
	Version   int                `json:"version,omitempty" graphql:"version,omitempty"`
	Data      map[string]Payload `json:"data,omitempty" graphql:"data,omitempty"`
}

type Data struct {
	Key     string  `json:"key"`
	Payload Payload `json:"payload"`
}

type Payload struct {

	//	Base 64 encoded value.
	Value interface{} `json:"value,omitempty"`

	//	Plaintext or Ciphertext
	Type Type `json:"type,omitempty"`
}

func (p *Payload) Map() map[string]interface{} {
	return map[string]interface{}{
		"value": p.Value,
		"type":  p.Type,
	}
}

func (d *Data) String() string {
	return fmt.Sprintf("%s=%v", d.Key, d.Payload.Value)
}

func (d *Data) GetPayload() *Payload {
	return &d.Payload
}

//	Returns KEY=VALUE mapping of the secret.
func (d *Data) KVMap() map[string]interface{} {
	return map[string]interface{}{
		d.Key: d.Payload.Value,
	}
}

//	Returns KEY=Payload{ Type, Value } mapping of the secret.
func (d *Data) Map() map[string]interface{} {
	return map[string]interface{}{
		d.Key: d.Payload,
	}
}

type Path struct {
	Organisation string `json:"org"`
	Project      string `json:"project"`
	Environment  string `json:"env"`
}

func (p *Path) Location() string {
	return fmt.Sprintf("%s/%s/%s", p.Organisation, p.Project, p.Environment)
}

type GenerateKeyOptions struct {

	//	Whether the key can be exported in the future.
	Exportable bool `json:"exportable,omitempty"`

	//	Key Type. We are using "aes256-gcm96" as the default one.
	Type string `json:"type,omitempty" default:"aes256-gcm96"`
}

func (o *GenerateKeyOptions) Marshal() ([]byte, error) {
	return json.Marshal(o)
}

type VaultResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Ciphertext string `json:"ciphertext,omitempty"`
		Plaintext  string `json:"plaintext,omitempty"`
		KeyVersion int    `json:"key_version"`
	} `json:"data"`
}

type SetRequestOptions struct {
	OrgID      string `json:"org_id"`
	EnvID      string `json:"env_id"`
	Data       Data   `json:"data"`
	KeyVersion int    `json:"key_version,omitempty"`
}

func (r *SetRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SetSecretOptions struct {
	KeyPath    string `json:"key_path"`
	EnvID      string `json:"env_id"`
	Data       Data   `json:"data"`
	KeyVersion int    `json:"key_version,omitempty"`
}

func (r *SetSecretOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (s *SetSecretOptions) GetVaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"plaintext":   s.Data.Payload.Value,
		"key_version": s.KeyVersion,
	}
}

type GetSecretOptions struct {
	Data    Data   `json:"data"`
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

type DecryptSecretOptions struct {
	Data        Data   `json:"data"`
	KeyLocation string `json:"key_location"`
	EnvID       string `json:"env_id"`
}

func (g *DecryptSecretOptions) GetVaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"ciphertext": g.Data.Payload.Value,
	}
}

type GetRequestOptions struct {
	Path    Path   `json:"path"`
	Key     string `json:"key"`
	Version *int   `json:"version,omitempty"`
}

func (r *GetRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetAllResponse struct {
	Data    map[string]Payload `json:"data"`
	Version int                `json:"version,omitempty"`
}

type ListRequestOptions struct {
	Path    Path `json:"path"`
	Version *int `json:"version,omitempty"`
}

func (r *ListRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
