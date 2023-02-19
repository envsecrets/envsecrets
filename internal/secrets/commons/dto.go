package commons

import (
	"encoding/json"
	"fmt"
)

type Secrets []Secret

func (s *Secrets) Map() map[string]interface{} {

	response := make(map[string]interface{})

	for _, item := range *s {
		response[item.Key] = item.Value
	}

	return response
}

type Secret struct {
	Key string `json:"key"`

	//	Base 64 encoded value.
	Value interface{} `json:"value"`
}

func (s *Secret) String() string {
	return fmt.Sprintf("%s=%v", s.Key, s.Value)
}

func (s *Secret) Map() map[string]interface{} {
	return map[string]interface{}{
		s.Key: s.Value,
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

type SetOptions struct {
	Path       Path   `json:"path"`
	Secret     Secret `json:"secret"`
	KeyVersion int    `json:"key_version"`
}

func (s *SetOptions) VaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"plaintext":   s.Secret.Value,
		"key_version": s.KeyVersion,
	}
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

type SetRequest struct {
	Path   Path   `json:"path"`
	Secret Secret `json:"secret"`
}

func (r *SetRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SetRequestOptions struct {
	EnvID  string `json:"env_id"`
	Secret Secret `json:"secret"`
}

type GetOptions struct {
	Secret  Secret `json:"secret"`
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

func (g *GetOptions) VaultOptions() map[string]interface{} {
	return map[string]interface{}{
		"ciphertext": g.Secret.Value,
	}
}

type GetRequestOptions struct {
	Path    Path   `json:"path"`
	Key     string `json:"key"`
	Version *int   `json:"version,omitempty"`
}

type GetAllResponse struct {
	Data    map[string]interface{} `json:"data"`
	Version int                    `json:"version,omitempty"`
}

type ListRequest struct {
	Path    Path `json:"path"`
	Version *int `json:"version,omitempty"`
}

func (r *ListRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
