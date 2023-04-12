package commons

import (
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

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

func MapToStruct(source any, target interface{}) error {
	if source == nil {
		return errors.New("source is nil")
	}
	entity, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(entity, &target)
}
