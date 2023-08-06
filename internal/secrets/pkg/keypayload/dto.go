package keypayload

import (
	"encoding/json"
	"fmt"

	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keyvalue"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/payload"
)

// Key-Payload Map
type KPMap map[string]*payload.Payload

// Returns a boolean validating whether the length of the map is 0.
func (m KPMap) IsEmpty() bool {
	return len(m) == 0
}

func (m KPMap) Load(value map[string]*payload.Payload) {
	for name := range value {
		m.Set(name, value[name])
	}
}

// Sets a key=payload pair to the map.
func (m KPMap) Set(key string, value *payload.Payload) {
	if m == nil {
		m = KPMap{}
	}
	m[key] = value
}

// Sets the value for the payload belong the the specified key in the map.
func (m KPMap) SetValue(key, value string) {
	m[key].Set(value)
}

// Fetches the payload for a specific key from the map.
func (m KPMap) Get(key string) *payload.Payload {
	return m[key]
}

// Fetches the value from the payload for a specific key in the map.
func (m KPMap) GetValue(key string) string {
	return m[key].GetValue()
}

// Returns string representation in the form of "key=value"
func (m KPMap) FmtString(key string) string {
	return fmt.Sprintf("%s=%s", key, m.GetValue(key))
}

// Updates a key name in the map from "old" to "new."
func (m KPMap) ChangeKey(old, new string) {
	payload := m.Get(old)
	m.Set(new, payload)
	m.Delete(old)
}

// Deletes a key=value pair from the map.
func (m KPMap) Delete(key string) {
	delete(m, key)
}

// Empties the values from the payloads of all key=value pairs.
func (m KPMap) DeleteValues() {
	for _, payload := range m {
		payload.DeleteValue()
	}
}

// Ovewrites or replaces values in the map for respective keys from supplied map.
func (m KPMap) Overwrite(source *KPMap) {
	for name, payload := range *source {
		m.Set(name, payload)
	}
}

// Base64 encodes all the pairs in the map.
func (m KPMap) Encode() {
	for name := range m {
		m[name].Encode()
	}
}

// Base64 decodes all the pairs in the map.
func (m KPMap) Decode() error {
	for name := range m {
		if err := m[name].Decode(); err != nil {
			return err
		}
	}
	return nil
}

// Encrypts all the key=value pairs with the provided key.
func (m KPMap) Encrypt(key [32]byte) error {
	for name := range m {
		if err := m[name].Encrypt(key); err != nil {
			return err
		}
	}
	return nil
}

// Encrypts all the key=value pairs with the provided key and returns a new deep copy of the map.
func (m KPMap) Encrypted(key [32]byte) (KPMap, error) {
	new := m
	if err := m.Encrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}

// Decrypts all the key=value pairs with the provided key.
func (m KPMap) Decrypt(key [32]byte) error {
	for name := range m {
		if err := m[name].Decrypt(key); err != nil {
			return err
		}
	}
	return nil
}

// Decrypts all the key=value pairs with the provided key and returns a new deep copy of the map.
func (m KPMap) Decrypted(key [32]byte) (KPMap, error) {
	new := m
	if err := m.Decrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}

// Returns a new key=value mapping.
func (m KPMap) ToKVMap() *keyvalue.KVMap {
	var result keyvalue.KVMap
	for name, payload := range m {
		result.Set(name, payload.Value)
	}
	return &result
}

// Returns JSON Marshalled map.
func (m *KPMap) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Marks payload value for specified key as Base64 encoded.
func (m KPMap) MarkEncoded(key string) {
	m[key].MarkEncoded()
}

// Marks payload value for specified key as Base64 decoded.
func (m KPMap) MarkDecoded(key string) {
	m[key].MarkDecoded()
}

// Marks all payload values as Base64 encoded.
func (m KPMap) MarkAllEncoded() {
	for _, payload := range m {
		payload.MarkEncoded()
	}
}

// Marks all payload values as Base64 decoded.
func (m KPMap) MarkAllDecoded() {
	for _, payload := range m {
		payload.MarkDecoded()
	}
}

// Marks the "exposable" value of the payload as "true."
//
//	Exposability allows the value to be synced as an exposable one
//	on platforms which differentiate between decryptable and non-decryptable secrets.
//	For example, Github and Vercel.
//	In Github actions, this value will be synced as a "variable" and NOT a secret, once it is marked "exposable" over here.
func (m KPMap) MarkExposable(key string) {
	m[key].MarkExposable()
}

func (m KPMap) MarkNotExposable(key string) {
	m[key].MarkNotExposable()
}
