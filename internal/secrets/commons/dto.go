package commons

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/envsecrets/envsecrets/internal/secrets/internal/keypayload"
	"github.com/envsecrets/envsecrets/internal/secrets/internal/keyvalue"
	"github.com/envsecrets/envsecrets/internal/secrets/internal/payload"
)

type Secret struct {

	//	To allows mutually exclusive locking
	sync.Mutex

	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	EnvID     string    `json:"env_id,omitempty"`

	//	Version of the secrets, if available.
	Version *int `json:"version,omitempty"`

	// Contains the secret mapping.
	Data keypayload.KPMap `json:"data,omitempty"`
}

// Returns the secret's key=value mapping.
func (s *Secret) GetMap() map[string]*payload.Payload {
	return s.Data
}

// Checks whether the secret contains even a single key=value mapping.
func (s *Secret) IsEmpty() bool {
	return len(s.Data) == 0
}

// Returns a new initialized 'Secret' object.
func New() *Secret {
	return &Secret{}
}

func ParseAndInitialize(data []byte) (*Secret, error) {

	var result Secret
	if data == nil {
		return nil, fmt.Errorf("invalid inputs")
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	result.MarkEncoded()
	return &result, nil
}

// Sets a key=value pair to the map.
func (s *Secret) Set(key string, value *payload.Payload) {
	if s.Data == nil {
		s.Data = make(keypayload.KPMap)
	}
	s.Data.Set(key, value)
}

type AddConfig struct {
	Value string `json:"value,omitempty"`

	//	Allows the value to be synced as an exposable one
	//	on platforms which differentiate between decryptable and non-decryptable secrets.
	//	For example, Github and Vercel.
	Exposable bool `json:"exposable,omitempty"`
}

// Sets a new key=value pair to the map.
func (s *Secret) Add(key string, config *AddConfig) {
	s.Set(key, &payload.Payload{
		Value:     config.Value,
		Exposable: config.Exposable,
	})
}

// Fetches the value for a specific key from the map.
func (s *Secret) Get(key string) *payload.Payload {
	return s.Data.Get(key)
}

// Fetches the value for a specific key from the map.
func (s *Secret) GetValue(key string) string {
	return s.Data.GetValue(key)
}

//	Get formatted string.
//
// Fetches key=value representation for a specific key and value from the map.
func (s *Secret) GetFmtString(key string) string {
	return fmt.Sprintf("%s=%s", key, s.Data[key].Value)
}

// Deletes a key=value pair from the map.
func (s *Secret) Delete(key string) {
	delete(s.Data, key)
}

// Ovewrites or replaces values in the map for respective keys from supplied map.
func (s *Secret) Overwrite(source map[string]*payload.Payload) {
	m := make(keypayload.KPMap)
	for name, payload := range source {
		m.Set(name, payload)
	}
	s.Data.Overwrite(&m)
}

// Base64 encodes all the pairs in the map.
func (s *Secret) Encode() {
	s.Data.Encode()
}

// Base64 decodes all the pairs in the map.
func (s *Secret) Decode() error {
	return s.Data.Decode()
}

// Marks all payload values as Base64 encoded.
func (s *Secret) MarkEncoded() {
	s.Data.MarkEncoded()
}

// Encrypts all the key=value pairs with provided encryption key.
func (s *Secret) Encrypt(key [32]byte) error {
	return s.Data.Encrypt(key)
}

// Encrypts all the key=value pairs with provided key
// and returns a new deep copy of the secret data without mutating the existing one.
func (s *Secret) Encrypted(key [32]byte) (*Secret, error) {
	new := s
	if err := new.Encrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}

// Decrypts all the key=value pairs with provided decryption key.
func (s *Secret) Decrypt(key [32]byte) error {
	return s.Data.Decrypt(key)
}

// Decrypts all the key=value pairs with provided key
// and returns a new deep copy of the secret data without mutating the existing one.
func (s *Secret) Decrypted(key [32]byte) (result *Secret, err error) {
	new := s
	if err := new.Decrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}

// Empties the values from the payloads of all key=value pairs.
func (s *Secret) Empty() {
	s.Lock()
	defer s.Unlock()

	for _, payload := range s.Data {
		payload.Empty()
	}
}

// Converts all the key=value pairs to a map.
func (s *Secret) ToMap() *keyvalue.KVMap {
	var result keyvalue.KVMap
	for name, payload := range s.Data {
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
	EnvID      string                      `json:"env_id"`
	Data       map[string]*payload.Payload `json:"data"`
	KeyVersion *int                        `json:"key_version,omitempty"`
}

func (r *SetRequestOptions) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SetSecretOptions struct {
	EnvID      string                      `json:"env_id"`
	Data       map[string]*payload.Payload `json:"data"`
	KeyVersion *int                        `json:"key_version,omitempty"`
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
	Secret Secret `json:"secret"`
	OrgID  string `json:"org"`
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
	Secret  Secret `json:"secret"`
	Version *int   `json:"version,omitempty"`
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
