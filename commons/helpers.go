package commons

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

//	Get the JWT key from NHOST variables
func GetJWTSecret() (*JWTSecret, error) {
	var response JWTSecret
	payload := os.Getenv("NHOST_JWT_SECRET")
	if err := json.Unmarshal([]byte(payload), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func GetAESKey() (string, error) {
	payload := os.Getenv("AES_KEY")
	if payload == "" {
		return "", errors.New("invalid key")
	}
	return payload, nil
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidateHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MapToStruct(source interface{}, target interface{}) error {
	if source == nil {
		return errors.New("source is nil")
	}
	entity, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(entity, &target)
}
