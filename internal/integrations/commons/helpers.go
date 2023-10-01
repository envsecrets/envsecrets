package commons

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
)

func EncryptCredentials(ctx context.ServiceContext, org_id string, payload map[string]interface{}) ([]byte, error) {

	//	Prepare the credentials
	credentials, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	//	Use server's own key to encrypt the credentials.
	kBytes, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PUBLIC_KEY"))
	if err != nil {
		return nil, err
	}

	var k [32]byte
	copy(k[:], kBytes)

	//	Encrypt the secrets with the server-copy of organisation's encryption key.
	encrypted, err := keys.SealAsymmetricallyAnonymous(credentials, k)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func DecryptCredentials(ctx context.ServiceContext, org_id string, payload []byte) ([]byte, error) {

	//	Get the server's copy of organisation's encryption key.
	kBytes, err := keys.GetOrgKeyServerCopy(ctx, org_id)
	if err == nil {
		var k [32]byte
		copy(k[:], kBytes)
		return keys.OpenSymmetrically(payload, k)
	}

	//	If the credentials were encrypted with server's public key,
	//	us thee server's private key to decrypt the credentials.
	var privateKey, publicKey [32]byte
	privateKeyBytes, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}
	copy(privateKey[:], privateKeyBytes)

	publicKeyBytes, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PUBLIC_KEY"))
	if err != nil {
		return nil, err
	}
	copy(publicKey[:], publicKeyBytes)

	//	Decrypt the value using org-key.
	return keys.OpenAsymmetricallyAnonymous(payload, publicKey, privateKey)
}
