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

type VaultResponse struct {
	Errors        []interface{} `json:"errors"`
	RequestID     string        `json:"request_id,omitempty"`
	LeaseID       string        `json:"lease_id,omitempty"`
	Renewable     bool          `json:"renewable,omitempty"`
	LeaseDuration int           `json:"lease_duration,omitempty"`
	Data          struct {
		Ciphertext string `json:"ciphertext,omitempty"`
		Plaintext  string `json:"plaintext,omitempty"`
		KeyVersion int    `json:"key_version,omitempty"`
		Backup     string `json:"backup,omitempty"`
	} `json:"data,omitempty"`
}

type SetRequestOptions struct {
	EnvID      string             `json:"env_id"`
	Data       map[string]Payload `json:"data"`
	KeyVersion *int               `json:"key_version,omitempty"`
}

func (r *SetRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SetSecretOptions struct {
	KeyPath    string             `json:"key_path"`
	EnvID      string             `json:"env_id"`
	Data       map[string]Payload `json:"data"`
	KeyVersion *int               `json:"key_version,omitempty"`
}

func (r *SetSecretOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (s *SetSecretOptions) GetVaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"plaintext":   s.Data,
		"key_version": s.KeyVersion,
	}
}

type CleanupSecretOptions struct {
	EnvID   string `json:"env_id"`
	Version int    `json:"version"`
}

type DeleteSecretOptions struct {
	EnvID   string `json:"env_id"`
	Key     string `json:"key"`
	Version *int   `json:"version"`
}

type DeleteRequestOptions struct {
	EnvID   string `query:"env_id"`
	Key     string `json:"key"`
	Version *int   `json:"version"`
}

func (r *DeleteRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type DecryptSecretOptions struct {
	Value       interface{} `json:"value"`
	KeyLocation string      `json:"key_location"`
}

func (g *DecryptSecretOptions) GetVaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"ciphertext": g.Value,
	}
}

type GetRequestOptions struct {
	EnvID   string `query:"env_id"`
	Key     string `query:"key"`
	Version *int   `query:"version"`
}

func (r *GetRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetSecretOptions struct {
	Key     string `json:"key"`
	KeyPath string `json:"key_path"`
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

type GetResponse struct {
	Data    map[string]Payload `json:"data"`
	Version *int               `json:"version,omitempty"`
}

type MergeRequestOptions struct {
	SourceEnvID   string `json:"source_env_id"`
	SourceVersion *int   `json:"source_version"`
	TargetEnvID   string `query:"env_id"`
}

func (r *MergeRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MergeSecretOptions struct {
	SourceEnvID   string `json:"source_env_id"`
	SourceVersion *int   `json:"source_version"`
	TargetEnvID   string `json:"target_env_id"`
}

type MergeResponse struct {
	Version *int `json:"version,omitempty"`
}

type ListRequestOptions struct {
	EnvID   string `query:"env_id"`
	Version *int   `query:"version,omitempty"`
}

func (r *ListRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
