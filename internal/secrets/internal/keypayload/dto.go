package keypayload

/*
import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/secrets/internal/keyvalue"
	"github.com/envsecrets/envsecrets/internal/secrets/internal/payload"
)

//	Key-Payload Map
//
// Contains the following mapping:
// key : { value: secret-value, type: plaintext/ciphertext }
type KPMap map[string]*payload.Payload

// Sets a key=payload pair to the map.
func (m KPMap) Set(key string, value payload.Payload) {
	m[key] = value
}

// Fetches the payload for a specific key from the map.
func (m KPMap) Get(key string) payload.Payload {
	return m[key]
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

// Encrypts all the key=value pairs with provided decryption key.
func (m KPMap) Encrypt(key [32]byte) error {
	for name, payload := range m {
		if err := payload.Encrypt(key); err != nil {
			return err
		}
		m.Set(name, payload)
	}
	return nil
}

// Decrypts all the key=value pairs with provided decryption key.
func (m KPMap) Decrypt(key [32]byte) error {
	for name, payload := range m {
		if err := payload.Decrypt(key); err != nil {
			return err
		}
		m.Set(name, payload)
	}
	return nil
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
*/
