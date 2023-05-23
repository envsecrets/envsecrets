package commons

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/envsecrets/envsecrets/internal/keys"
)

type Type string

type Row struct {
	ID        string    `json:"id,omitempty" graphql:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" graphql:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" graphql:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty" graphql:"name,omitempty"`
	UserID    string    `json:"user_id,omitempty" graphql:"user_id,omitempty"`
	EnvID     string    `json:"env_id,omitempty" graphql:"env_id,omitempty"`
	Version   int       `json:"version,omitempty" graphql:"version,omitempty"`
	Data      Secrets   `json:"data,omitempty" graphql:"data,omitempty"`
}

// Contains the following mapping:
// key : { value: secret-value, type: plaintext/ciphertext }
type Secrets map[string]Payload

// Sets a key=value pair to the map.
func (s Secrets) Set(key string, payload Payload) {
	s[key] = payload
}

// Fetches the value for a specific key from the map.
func (s Secrets) Get(key string) Payload {
	return s[key]
}

// Deletes a key=value pair from the map.
func (s Secrets) Detele(key string) {
	delete(s, key)
}

// Encrypts all the key=value pairs with provided encryption key.
func (s Secrets) Encrypt(key [32]byte) error {
	for name, payload := range s {
		if payload.Type == Ciphertext {
			encrypted, err := keys.SealSymmetrically([]byte(fmt.Sprintf("%v", payload.Value)), key)
			if err != nil {
				return err.Error
			}
			payload.Value = base64.StdEncoding.EncodeToString(encrypted)
		} else {
			payload.Value = base64.StdEncoding.EncodeToString([]byte(payload.Value))
		}
		s.Set(name, payload)
	}
	return nil
}

// Encrypts all the key=value pairs with provided encryption key
// and returns a new copy of the secrets payload without mutating the existing one.
func (s Secrets) Encrypted(key [32]byte) (result *Secrets, err error) {
	copy := s
	if err := copy.Encrypt(key); err != nil {
		return nil, err
	}
	return &copy, nil
}

// Decrypts all the key=value pairs with provided decryption key.
func (s Secrets) Decrypt(key [32]byte) error {
	for name, payload := range s {
		if payload.Type == Ciphertext {

			//	Base64 decode the secret value
			decoded, err := base64.StdEncoding.DecodeString(payload.Value)
			if err != nil {
				return err
			}

			//	Decrypt the value using org-key.
			decrypted, er := keys.OpenSymmetrically(decoded, key)
			if er != nil {
				return er.Error
			}

			payload.Value = base64.StdEncoding.EncodeToString(decrypted)
			s.Set(name, payload)
		}
	}
	return nil
}

// Decrypts all the key=value pairs with provided decryption key
// and returns a new copy of the secrets payload without mutating the existing one.
func (s Secrets) Decrypted(key [32]byte) (result *Secrets, err error) {
	copy := s
	if err := copy.Decrypt(key); err != nil {
		return nil, err
	}
	return &copy, nil
}

// Ovewrites or replaces values in the map for respective keys from supplied map.
func (s *Secrets) Overwrite(source Secrets) {
	for name, payload := range source {
		s.Set(name, payload)
	}
}

type Payload struct {

	//	Internal variable to record the current state of encoding of this payload's value.
	isEncoded bool `json:"-"`

	//	Value
	Value string `json:"value,omitempty"`

	//	Type in which the value is stored in our database.
	Type Type `json:"type,omitempty"`
}

func (p *Payload) GetDecodedValue() ([]byte, error) {
	return base64.StdEncoding.DecodeString(p.Value)
}

// Converts the payload to a map.
func (p *Payload) Map() (result map[string]interface{}, err error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(payload, &result)
	return
}

// Map representing key=value pairs.
type KVMap map[string]string

// Sets a key=value pair to the map.
func (m KVMap) Set(key, value string) {
	m[key] = value
}

// Fetches the value for a specific key from the map.
func (m KVMap) Get(key string) string {
	return m[key]
}

// Deletes a key=value pair from the map.
func (m KVMap) Detele(key string) {
	delete(m, key)
}

// Returns JSON Marshalled map.
func (m *KVMap) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshalls the json in provided interface.
func (m *KVMap) UnmarshalIn(i interface{}) error {
	payload, err := m.Marshal()
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, &i)
}

// Converts all the key=value pairs to a map.
func (s *Secrets) ToMap() *KVMap {
	result := make(KVMap)
	for name, payload := range *s {
		result.Set(name, payload.Value)
	}
	return &result
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
	EnvID      string  `json:"env_id"`
	Secrets    Secrets `json:"secrets"`
	KeyVersion *int    `json:"key_version,omitempty"`
}

func (r *SetRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SetSecretOptions struct {
	EnvID      string  `json:"env_id"`
	Secrets    Secrets `json:"secrets"`
	KeyVersion *int    `json:"key_version,omitempty"`
}

func (r *SetSecretOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
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
	Secrets Secrets `json:"secrets"`
	OrgID   string  `json:"org"`
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
	EnvID   string `json:"env_id"`
	Version *int   `json:"version,omitempty"`
}

type GetResponse struct {
	Secrets Secrets `json:"secrets"`
	Version *int    `json:"version,omitempty"`
}

type MergeRequestOptions struct {
	SourceEnvID   string `json:"source_env_id"`
	SourceVersion *int   `json:"source_version"`
	TargetEnvID   string `json:"env_id"`
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
