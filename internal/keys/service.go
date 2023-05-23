package keys

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"

	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/keys/graphql"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateOptions) error {
	return graphql.Create(ctx, client, options)
}

func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateWithUserIDOptions) error {
	return graphql.CreateWithUserID(ctx, client, options)
}

func GetByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) (*commons.Key, error) {
	return graphql.GetByUserID(ctx, client, user_id)
}

func GetPublicKeyByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) ([]byte, error) {
	return graphql.GetPublicKeyByUserID(ctx, client, user_id)
}

func SealSymmetrically(message []byte, key [commons.KEY_BYTES]byte) ([]byte, error) {

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [commons.NONCE_LEN]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	return secretbox.Seal(nonce[:], message, &nonce, &key), nil
}

func OpenSymmetrically(message []byte, key [commons.KEY_BYTES]byte) ([]byte, error) {

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [commons.NONCE_LEN]byte
	copy(nonce[:], message[:commons.NONCE_LEN])

	result, ok := secretbox.Open(nil, message[commons.NONCE_LEN:], &nonce, &key)
	if !ok {
		return nil, errors.New("failed to open the message from symmetric key")
	}

	return result, nil
}

func SealAsymmetricallyAnonymous(message []byte, key [commons.KEY_BYTES]byte) ([]byte, error) {

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [commons.NONCE_LEN]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	// This encrypts msg and appends the result to the nonce.
	result, err := box.SealAnonymous(nonce[:], message, &key, rand.Reader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func OpenAsymmetricallyAnonymous(message []byte, publicKey, privateKey [commons.KEY_BYTES]byte) ([]byte, error) {

	// The recipient can decrypt the message using their private key and the
	// sender's public key. When you decrypt, you must use the same nonce you
	// used to encrypt the message. One way to achieve this is to store the
	// nonce alongside the encrypted message. Above, we stored the nonce in the
	// first 24 bytes of the encrypted text.
	var nonce [commons.NONCE_LEN]byte
	copy(nonce[:], message[:commons.NONCE_LEN])
	result, ok := box.OpenAnonymous(nil, message[commons.NONCE_LEN:], &publicKey, &privateKey)
	if !ok {
		return nil, errors.New("failed to open the message from asymmetric key")
	}

	return result, nil
}

func GenerateKeyPair(password string) (*commons.IssueKeyPairResponse, error) {

	publicKeyBytes, privateKeyBytes, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	//	Generate a separate random symmetric key
	protectionKeyBytes, err := globalCommons.GenerateRandomBytes(commons.KEY_BYTES)
	if err != nil {
		return nil, err
	}

	//	Encrypt the private key using protection key
	var protectionKeyForSealing [32]byte
	copy(protectionKeyForSealing[:], protectionKeyBytes)
	encryptedPrivateKeyBytes, err := SealSymmetrically(privateKeyBytes[:], protectionKeyForSealing)
	if err != nil {
		return nil, err
	}

	//	Generate random 32 byte salt
	saltBytes, err := globalCommons.GenerateRandomBytes(commons.KEY_BYTES)
	if err != nil {
		return nil, err
	}

	//	Safegaurd the protection key using Argon2i password based hashing.
	passwordDerivedKey := argon2.Key([]byte(password), saltBytes, 3, commons.KEY_BYTES*1024, 4, commons.KEY_BYTES)

	//	Encrypt the protection key using password derived key
	var passwordDerivedKeyForSealing [32]byte
	copy(passwordDerivedKeyForSealing[:], passwordDerivedKey)
	encryptedProtectionKeyBytes, err := SealSymmetrically(protectionKeyBytes, passwordDerivedKeyForSealing)
	if err != nil {
		return nil, err
	}

	return &commons.IssueKeyPairResponse{
		PublicKey:           publicKeyBytes[:],
		PrivateKey:          encryptedPrivateKeyBytes,
		DecryptedPrivateKey: privateKeyBytes[:],
		ProtectedKey:        encryptedProtectionKeyBytes,
		Salt:                saltBytes,
	}, nil
}

// Decrypt the org's symmetric key with your local public-private key.
func DecryptAsymmetricallyAnonymous(public, private, org_key []byte) ([]byte, error) {

	var publicKey, privateKey [commons.KEY_BYTES]byte
	copy(publicKey[:], public)
	copy(privateKey[:], private)
	result, err := OpenAsymmetricallyAnonymous(org_key, publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetOrgKeyServerCopy(ctx context.ServiceContext, org_id string) ([]byte, error) {

	//	Initialize new GQL client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Get the server's key copy
	serverCopy, err := organisations.GetServerKeyCopy(ctx, client, org_id)
	if err != nil {
		return nil, err
	}

	//	Decrypt the copy with server's private key (in env vars).
	serverPrivateKey, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}

	serverPublicKey, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PUBLIC_KEY"))
	if err != nil {
		return nil, err
	}

	return DecryptAsymmetricallyAnonymous(serverPublicKey, serverPrivateKey, serverCopy)
}
