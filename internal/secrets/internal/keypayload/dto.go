package keypayload

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/secrets/internal/keyvalue"
	"github.com/envsecrets/envsecrets/internal/secrets/internal/payload"
)

// Key-Payload Map
type KPMap map[string]*payload.Payload

// Sets a key=payload pair to the map.
func (m KPMap) Set(key string, value *payload.Payload) {
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

// Deletes a key=value pair from the map.
func (m KPMap) Detele(key string) {
	delete(m, key)
}

// Ovewrites or replaces values in the map for respective keys from supplied map.
func (m KPMap) Overwrite(source *KPMap) {
	for name, payload := range *source {
		m.Set(name, payload)
	}
}

// Base64 encodes all the pairs in the map.
func (m KPMap) Encode() {
	for name, payload := range m {
		payload.Encode()
		m.Set(name, payload)
	}
}

// Base64 decodes all the pairs in the map.
func (m KPMap) Decode() error {
	for name, payload := range m {
		if err := payload.Decode(); err != nil {
			return err
		}
		m.Set(name, payload)
	}
	return nil
}

// Encrypts all the key=value pairs with the provided key.
func (m KPMap) Encrypt(key [32]byte) error {
	for name, payload := range m {
		if err := payload.Encrypt(key); err != nil {
			return err
		}
		m.Set(name, payload)
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
	for name, payload := range m {
		if err := payload.Decrypt(key); err != nil {
			return err
		}
		m.Set(name, payload)
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

// Marks all payload values as Base64 encoded.
func (m KPMap) MarkEncoded() {
	for _, payload := range m {
		payload.MarkEncoded()
	}
}
