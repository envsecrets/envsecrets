package dto

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Map representing key=value pairs.
type KVMap struct {
	sync.Mutex

	mapping map[string]string
}

// Custom Unmarshaler.
func (m *KVMap) UnmarshalJSON(data []byte) error {

	var mapping map[string]string
	if err := json.Unmarshal(data, &mapping); err != nil {
		return err
	}

	*m = KVMap{
		mapping: mapping,
	}

	return nil
}

// Custom Marshaller.
func (m *KVMap) MarshalJSON() ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	return json.Marshal(m.mapping)
}

// Unmarshalls the json in provided interface.
func (m *KVMap) UnmarshalIn(i interface{}) error {
	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, &i)
}

// Sets a key=value pair to the map.
func (m *KVMap) Set(key, value string) {
	m.Lock()
	defer m.Unlock()
	if m.mapping == nil {
		m.mapping = make(map[string]string)
	}
	m.mapping[key] = value
}

// Fetches the value for a specific key from the map.
func (m *KVMap) Get(key string) string {
	m.Lock()
	defer m.Unlock()
	return m.mapping[key]
}

// Returns string representation in the form of "key=value"
func (m *KVMap) String(key string) string {
	m.Lock()
	defer m.Unlock()
	return fmt.Sprintf("%s=%s", key, m.mapping[key])
}

// Deletes a key=value pair from the map.
func (m *KVMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.mapping, key)
}

// Returns the mapping from the KVMap.
func (m *KVMap) GetMapping() map[string]string {
	m.Lock()
	defer m.Unlock()
	return m.mapping
}
