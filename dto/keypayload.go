package dto

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Key-Payload Map
type KPMap struct {
	sync.Mutex
	mapping map[string]*Payload
}

// Custom Unmarshaler.
func (m *KPMap) UnmarshalJSON(data []byte) error {

	var mapping map[string]*Payload
	if err := json.Unmarshal(data, &mapping); err != nil {
		return err
	}

	*m = KPMap{
		mapping: mapping,
	}

	return nil
}

// Custom Marshaller.
func (m *KPMap) MarshalJSON() ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	return json.Marshal(m.mapping)
}

// Returns an array of all the keys in the mapping.
func (m *KPMap) Keys() []string {
	m.Lock()
	defer m.Unlock()
	var keys []string
	for key := range m.mapping {
		keys = append(keys, key)
	}
	return keys
}

// Returns a boolean validating whether the length of the map is 0.
func (m *KPMap) IsEmpty() bool {
	return len(m.mapping) == 0
}

func (m *KPMap) Load(mapping *KPMap) {
	m.Lock()
	defer m.Unlock()
	for key, value := range mapping.mapping {
		m.Set(key, value)
	}
}

// Sets a key=payload pair to the map.
func (m *KPMap) Set(key string, value *Payload) {
	m.Lock()
	defer m.Unlock()
	if m.mapping == nil {
		m.mapping = make(map[string]*Payload)
	}
	m.mapping[key] = value
}

// Sets the value for the payload belong the the specified key in the map.
func (m *KPMap) SetValue(key, value string) {
	m.mapping[key].Set(value)
}

// Fetches the payload for a specific key from the map.
func (m *KPMap) Get(key string) *Payload {
	m.Lock()
	defer m.Unlock()
	return m.mapping[key]
}

// Fetches the value from the payload for a specific key in the map.
func (m *KPMap) GetValue(key string) string {
	return m.mapping[key].GetValue()
}

// Returns string representation in the form of "key=value"
func (m *KPMap) FmtString(key string) string {
	return fmt.Sprintf("%s=%s", key, m.GetValue(key))
}

// Returns an array of string representations in the form of "key=value"
func (m *KPMap) FmtStrings() []string {
	var result []string
	for key := range m.mapping {
		result = append(result, m.FmtString(key))
	}
	return result
}

// Updates a key name in the map from "old" to "new."
func (m *KPMap) ChangeKey(old, new string) {
	payload := m.Get(old)
	m.Set(new, payload)
	m.Delete(old)
}

// Deletes a key=value pair from the map.
func (m *KPMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.mapping, key)
}

// Empties the values from the payloads of all key=value pairs.
func (m *KPMap) DeleteValues() {
	for _, payload := range m.mapping {
		payload.DeleteValue()
	}
}

// Ovewrites or replaces values in the map for respective keys from supplied map.
func (m *KPMap) Overwrite(source *KPMap) {
	for name, payload := range source.mapping {
		m.Set(name, payload)
	}
}

// Base64 encodes all the pairs in the map.
func (m *KPMap) Encode() {
	for name := range m.mapping {
		m.mapping[name].Encode()
	}
}

// Base64 decodes all the pairs in the map.
func (m *KPMap) Decode() error {
	for name := range m.mapping {
		if err := m.mapping[name].Decode(); err != nil {
			return err
		}
	}
	return nil
}

// Encrypts all the key=value pairs with the provided key.
func (m *KPMap) Encrypt(key [32]byte) error {
	for name := range m.mapping {
		if err := m.mapping[name].Encrypt(key); err != nil {
			return err
		}
	}
	return nil
}

// Encrypts all the key=value pairs with the provided key and returns a new deep copy of the map.
func (m *KPMap) Encrypted(key [32]byte) (*KPMap, error) {
	new := m
	if err := new.Encrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}

// Decrypts all the key=value pairs with the provided key.
func (m *KPMap) Decrypt(key [32]byte) error {
	for name := range m.mapping {
		if err := m.mapping[name].Decrypt(key); err != nil {
			return err
		}
	}
	return nil
}

// Decrypts all the key=value pairs with the provided key and returns a new deep copy of the map.
func (m *KPMap) Decrypted(key [32]byte) (*KPMap, error) {
	new := m
	if err := new.Decrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}

// Returns a new key=value mapping.
func (m *KPMap) ToKVMap() *KVMap {
	var result KVMap
	for name, payload := range m.mapping {
		result.Set(name, payload.Value)
	}
	return &result
}

// Marks payload value for specified key as Base64 encoded.
func (m *KPMap) MarkEncoded(key string) {
	m.mapping[key].MarkEncoded()
}

// Marks payload value for specified key as Base64 decoded.
func (m *KPMap) MarkDecoded(key string) {
	m.mapping[key].MarkDecoded()
}

// Marks all payload values as Base64 encoded.
func (m *KPMap) MarkAllEncoded() {
	for key := range m.mapping {
		m.mapping[key].MarkEncoded()
	}
}

// Marks all payload values as Base64 decoded.
func (m *KPMap) MarkAllDecoded() {
	for key := range m.mapping {
		m.mapping[key].MarkDecoded()
	}
}

// Marks the "exposable" value of the payload as "true."
//
//	Exposability allows the value to be synced as an exposable one
//	on platforms which differentiate between decryptable and non-decryptable secrets.
//	For example, Github and Vercel.
//	In Github actions, this value will be synced as a "variable" and NOT a secret, once it is marked "exposable" over here.
func (m *KPMap) MarkExposable(key string) {
	m.mapping[key].MarkExposable()
}

func (m *KPMap) MarkNotExposable(key string) {
	m.mapping[key].MarkNotExposable()
}

// Returns the mapping.
func (m *KPMap) GetMapping() map[string]*Payload {
	return m.mapping
}
