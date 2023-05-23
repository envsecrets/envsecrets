package vercel

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/nacl/box"
)

func generateGithuAppJWT(file []byte) (string, error) {

	// expires in 60 minutes
	expiration := time.Now().Add(time.Second * 600)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    os.Getenv("GITHUB_APP_ID"),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiration),
		Subject:   "envsecrets temporary github auth jwt",
	})

	pkey, err := jwt.ParseRSAPrivateKeyFromPEM(file)
	if err != nil {
		return "", err
	}

	response, err := token.SignedString(pkey)
	if err != nil {
		return "", err
	}

	return response, nil
}

// Encrypt a secret value using libsodium equivalent NACL secret box method.
func encryptSecret(pk, secret string) (string, error) {
	var pkBytes [32]byte
	copy(pkBytes[:], pk)
	secretBytes := []byte(secret)

	out := make([]byte, 0,
		len(secretBytes)+
			box.Overhead+
			len(pkBytes))

	enc, err := box.SealAnonymous(
		out, secretBytes, &pkBytes, rand.Reader,
	)
	if err != nil {
		return "", err
	}

	encEnc := base64.StdEncoding.EncodeToString(enc)
	return encEnc, nil
}
