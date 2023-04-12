package commons

import (
	"encoding/json"
)

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

type CleanupKeyOptions struct {
	EnvID   string `json:"env_id"`
	Version int    `json:"version"`
}

type DeleteKeyOptions struct {
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

type GetKeyOptions struct {
	Key     string `json:"key"`
	KeyPath string `json:"key_path"`
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

type MergeKeyOptions struct {
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

type KeyExportOptions struct {
	Type    string `json:"key_type"`
	Version string `json:"version"`
	Name    string `json:"name"`
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

type KeyExportResponse struct {
	Data struct {
		Name string                 `json:"name"`
		Keys map[string]interface{} `json:"keys"`
	} `json:"data"`
}
