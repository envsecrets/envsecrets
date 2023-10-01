package dto

import (
	"encoding/json"
	"fmt"
)

type Secret struct {
	Base

	//	The UUID of the user in our database who created this secret.
	//
	//	reference: https://en.wikipedia.org/wiki/Universally_unique_identifier
	//	required: true
	UserID string `json:"user_id,omitempty"`

	//	The UUID of the project environment in our database this secret belongs to.
	//
	//	reference: https://en.wikipedia.org/wiki/Universally_unique_identifier
	//	required: true
	EnvID string `json:"env_id,omitempty"`

	//	The version of this secret.
	//	When the secret is created, the default version is 1.
	//	While fetching the secret, sometimes the version may not be fetched, therefore, it is a pointer which can be nil.
	//
	//	required: false
	//	default: 0
	Version *int `json:"version,omitempty"`

	//	The mapping of all the key-value pairs this secret contains.
	//
	//	format: dto.KPMap
	//	required: true
	Data *KPMap `json:"data,omitempty"`
}

func (s *Secret) UnmarshalJSON(data []byte) error {

	var base Base
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	var result Secret

	var structure struct {
		UserID  string `json:"user_id,omitempty"`
		EnvID   string `json:"env_id,omitempty"`
		Version *int   `json:"version,omitempty"`
		Data    *KPMap `json:"data,omitempty"`
	}

	if err := json.Unmarshal(data, &structure); err != nil {
		return err
	}

	result.Base = base
	result.UserID = structure.UserID
	result.EnvID = structure.EnvID
	result.Version = structure.Version
	result.Data = structure.Data

	*s = result
	return nil
}

// Returns a shallow copy of the secret's key=value mapping.
func (s *Secret) DataCopy() KPMap {
	return KPMap{
		mapping: s.Data.mapping,
	}
}

// Checks whether the secret contains even a single key=value mapping.
func (s *Secret) IsEmpty() bool {
	return s.Data == nil || s.Data.IsEmpty()
}

// Sets a key=value pair to the map.
func (s *Secret) Set(key string, value *Payload) {
	if s.Data == nil {
		s.Data = &KPMap{}
	}
	s.Data.Set(key, value)
}

// Increases the version of the secret by 1.
func (s *Secret) IncrementVersion() {
	if s.Version != nil {
		*s.Version += 1
	}
}

// Fetches the value for a specific key from the map.
func (s *Secret) Get(key string) *Payload {
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
	return fmt.Sprintf("%s=%s", key, s.Data.mapping[key].Value)
}

// Deletes a key=value pair from the map.
func (s *Secret) Delete(key string) {
	s.Data.Delete(key)
}

// Updates a key name in the map from "old" to "new."
func (s *Secret) ChangeKey(old, new string) {
	s.Data.ChangeKey(old, new)
}

// Ovewrites or replaces values in the map for respective keys from supplied map.
func (s *Secret) Overwrite(source *KPMap) {
	m := KPMap{}
	m.Load(source)
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
	s.Data.MarkAllEncoded()
}

// Marks the "exposable" value of the payload as "true."
//
//	Exposability allows the value to be synced as an exposable one
//	on platforms which differentiate between decryptable and non-decryptable secrets.
//	For example, Github and Vercel.
//	In Github actions, this value will be synced as a "variable" and NOT a secret, once it is marked "exposable" over here.
func (s *Secret) MarkExposable(key string) {
	s.Data.MarkExposable(key)
}

func (s *Secret) MarkNotExposable(key string) {
	s.Data.MarkExposable(key)
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
func (s *Secret) DeleteValues() {
	s.Data.DeleteValues()
}

// Converts all the key=value pairs to a map.
func (s *Secret) ToKVMap() *KVMap {
	var result KVMap
	for name, payload := range s.Data.mapping {
		result.Set(name, payload.Value)
	}
	return &result
}
