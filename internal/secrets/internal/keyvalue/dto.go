package keyvalue

import (
	"encoding/json"
	"fmt"
)

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

// Returns string representation in the form of "key=value"
func (m KVMap) String(key string) string {
	return fmt.Sprintf("%s=%s", key, m[key])
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
