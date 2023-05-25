package payload

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	"github.com/envsecrets/envsecrets/internal/keys"
)

// Type of secret value.
// For example, "plaintext" or "ciphertext".
//type Type string

type Payload struct {
	sync.Mutex
	Value string `json:"value,omitempty"`

	//	Allows the value to be synced as an exposable one
	//	on platforms which differentiate between decryptable and non-decryptable secrets.
	//	For example, Github and Vercel.
	Exposable bool `json:"exposable,omitempty"`

	//	Type in which the value is stored in our database.
	//	Type Type `json:"type,omitempty"`

	//	Internal variable to record the current state of encoding of this payload's value.
	encoded bool `json:"-"`
}

// Converts the payload to a map.
func (p *Payload) ToMap() (result map[string]interface{}, err error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(payload, &result)
	return
}

// Sets the value of the payload.
func (p *Payload) Set(value string) {
	p.Lock()
	defer p.Unlock()
	p.Value = value
}

// Returns the value of the payload.
func (p *Payload) GetValue() string {
	p.Lock()
	defer p.Unlock()
	return p.Value
}

// Returns a boolean validating whether a value is exposable or not.
func (p *Payload) IsExposable() bool {
	return p.Exposable
}

// Marks the payload as "exportable."
func (p *Payload) MarkEncoded() {
	p.Lock()
	defer p.Unlock()
	p.encoded = true
}

// Empties the values from the payload.
func (p *Payload) DeleteValue() {
	p.Set("")
}

// Returns boolean indicating whether the value of the payload is empty or not.
func (p *Payload) IsEmpty() bool {
	return p.Value == ""
}

// Marks the "exposable" value of the payload as "true."
//
//	Exposability allows the value to be synced as an exposable one
//	on platforms which differentiate between decryptable and non-decryptable secrets.
//	For example, Github and Vercel.
//	In Github actions, this value will be synced as a "variable" and NOT a secret, once it is marked "exposable" over here.
func (p *Payload) MarkExposable() {
	p.Lock()
	defer p.Unlock()
	p.Exposable = true
}

// Marks the "exposable" value of the payload as "false."
//
//	Read the documentation of "MarkExposable" function.
func (p *Payload) MarkNotExposable() {
	p.Lock()
	defer p.Unlock()
	p.Exposable = false
}

// Marks the payload as "decoded."
func (p *Payload) MarkDecoded() {
	p.Lock()
	defer p.Unlock()
	p.encoded = false
}

// Returns boolean whether the payload is already encoded or not.
func (p *Payload) IsEncoded() bool {
	return p.encoded
}

// Base64 encodes the value of the payload.
func (p *Payload) Encode() {
	value := base64.StdEncoding.EncodeToString([]byte(p.Value))

	//	Update the value.
	p.Set(value)

	//	Mark the payload as encoded.
	p.MarkEncoded()
}

// Base64 decodes the value of the payload.
func (p *Payload) Decode() error {
	value, err := base64.StdEncoding.DecodeString(p.Value)
	if err != nil {
		return err
	}

	//	Update the value.
	p.Set(string(value))

	//	Mark the payload as decoded.
	p.MarkDecoded()
	return nil
}

func (p *Payload) Encrypt(key [32]byte) error {

	//	Decode the value before encrypting it.
	if p.encoded {
		if err := p.Decode(); err != nil {
			return err
		}
	}

	encrypted, err := keys.SealSymmetrically([]byte(p.Value), key)
	if err != nil {
		return err
	}
	p.Set(base64.StdEncoding.EncodeToString(encrypted))

	p.MarkEncoded()
	return nil
}

// Encrypts the value of the payload with the supplied key and returns a new deep copy of the payload.
func (p *Payload) Encrypted(key [32]byte) (*Payload, error) {
	new := p
	if err := new.Encrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}

// Decrypts the value of the payload with the supplied key.
func (p *Payload) Decrypt(key [32]byte) error {

	//	Decode the value before decrypting it.
	if p.encoded {
		if err := p.Decode(); err != nil {
			return err
		}
	}

	//	Decrypt the value using org-key.
	decrypted, err := keys.OpenSymmetrically([]byte(p.Value), key)
	if err != nil {
		return err
	}

	p.Set(string(decrypted))

	//	Encode the value once again after saving it.
	p.Encode()

	return nil
}

// Decrypts the value of the payload and returns a new deep copy of the payload.
func (p *Payload) Decrypted(key [32]byte) (*Payload, error) {
	new := p
	if err := new.Decrypt(key); err != nil {
		return nil, err
	}
	return new, nil
}
