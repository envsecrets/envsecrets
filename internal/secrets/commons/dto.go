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

	//	If set, enables taking backup of named key in the plaintext format. Once set, this cannot be disabled.
	AllowPlaintextBackup bool `json:"allow_plaintext_backup,omitempty"`

	//	Key Type. We are using "aes256-gcm96" as the default one.
	Type string `json:"type,omitempty" default:"aes256-gcm96"`
}

func (o *GenerateKeyOptions) Marshal() ([]byte, error) {
	return json.Marshal(o)
}

type VaultResponse struct {
	Errors        []interface{} `json:"errors,omitempty"`
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
	OrgID      string             `json:"org_id"`
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
	EnvID   string `json:"env_id"`
	Key     string `json:"key"`
	Version *int   `json:"version"`
}

func (r *DeleteRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
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
	OrgID   string `query:"org_id"`
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
	OrgID         string `json:"org_id"`
	SourceEnvID   string `json:"source_env_id"`
	SourceVersion *int   `json:"source_version"`
	TargetEnvID   string `json:"target_env_id"`
}

func (r *MergeRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MergeSecretOptions struct {
	KeyPath       string `json:"key_path"`
	SourceEnvID   string `json:"source_env_id"`
	SourceVersion *int   `json:"source_version"`
	TargetEnvID   string `json:"target_env_id"`
}

type MergeResponse struct {
	Version *int `json:"version,omitempty"`
}

type KeyRestoreRequestOptions struct {
	OrgID  string `json:"org_id"`
	Backup string `json:"backup"`
}

type KeyRestoreOptions struct {
	Backup string `json:"backup"`
}

func (r *KeyRestoreOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type VaultKeyExportResponse struct {
	Data struct {
		Name string                 `json:"name"`
		Keys map[string]interface{} `json:"keys"`
	} `json:"data"`
}

type KeyBackupRequestOptions struct {
	OrgID string `query:"org_id"`
}

type KeyConfigUpdateOptions struct {

	//	Specifies if the key is allowed to be deleted.
	DeletionAllowed bool `json:"deletion_allowed,omitempty"`

	//	Whether the key can be exported in the future.
	Exportable bool `json:"exportable,omitempty"`

	//	If set, enables taking backup of named key in the plaintext format. Once set, this cannot be disabled.
	AllowPlaintextBackup bool `json:"allow_plaintext_backup,omitempty"`
}

func (r *KeyConfigUpdateOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type KeyBackupResponse struct {
	Data struct {
		Backup string `json:"backup"`
	} `json:"data"`
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
