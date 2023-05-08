package keys

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	internalErrors "errors"
	"io"
	"os"

	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/keys/graphql"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateOptions) *errors.Error {
	return graphql.Create(ctx, client, options)
}

func CreateWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateWithUserIDOptions) *errors.Error {
	return graphql.CreateWithUserID(ctx, client, options)
}

func GetByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) (*commons.Key, *errors.Error) {
	return graphql.GetByUserID(ctx, client, user_id)
}

func GetPublicKeyByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) ([]byte, *errors.Error) {
	return graphql.GetPublicKeyByUserID(ctx, client, user_id)
}

func SealSymmetrically(message []byte, key [commons.KEY_BYTES]byte) []byte {

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [commons.NONCE_LEN]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	return secretbox.Seal(nonce[:], message, &nonce, &key)
}

func OpenSymmetrically(message []byte, key [commons.KEY_BYTES]byte) ([]byte, *errors.Error) {

	errMessage := "Failed to open the message from symmetric key"

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [commons.NONCE_LEN]byte
	copy(nonce[:], message[:commons.NONCE_LEN])

	result, ok := secretbox.Open(nil, message[commons.NONCE_LEN:], &nonce, &key)
	if !ok {
		return nil, errors.New(internalErrors.New(errMessage), errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return result, nil
}

func SealAsymmetricallyAnonymous(message []byte, key [commons.KEY_BYTES]byte) ([]byte, *errors.Error) {

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [commons.NONCE_LEN]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// This encrypts msg and appends the result to the nonce.
	result, err := box.SealAnonymous(nonce[:], message, &key, rand.Reader)
	if err != nil {
		return nil, errors.New(err, "Failed to seal the message anonymously", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return result, nil
}

func OpenAsymmetricallyAnonymous(message []byte, publicKey, privateKey [commons.KEY_BYTES]byte) ([]byte, *errors.Error) {

	errMessage := "Failed to open the message from asymmetric keys"
	// The recipient can decrypt the message using their private key and the
	// sender's public key. When you decrypt, you must use the same nonce you
	// used to encrypt the message. One way to achieve this is to store the
	// nonce alongside the encrypted message. Above, we stored the nonce in the
	// first 24 bytes of the encrypted text.
	var nonce [commons.NONCE_LEN]byte
	copy(nonce[:], message[:commons.NONCE_LEN])
	result, ok := box.OpenAnonymous(nil, message[commons.NONCE_LEN:], &publicKey, &privateKey)
	if !ok {
		return nil, errors.New(internalErrors.New(errMessage), errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return result, nil
}

func VerifyKeyPair(private, public []byte) (bool, error) {

	key, err := x509.ParsePKCS1PrivateKey(private)
	if err != nil {
		return false, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return false, err
	}
	return key.PublicKey.Equal(pubKey), nil
}

func GenerateKeyPair(password string) (*commons.IssueKeyPairResponse, *errors.Error) {

	publicKeyBytes, privateKeyBytes, er := box.GenerateKey(rand.Reader)
	if er != nil {
		return nil, errors.New(er, "Failed to generate public-private key pair", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	//	Generate a separate random symmetric key
	protectionKeyBytes, er := globalCommons.GenerateRandomBytes(commons.KEY_BYTES)
	if er != nil {
		return nil, errors.New(er, "Failed to generate protection key", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	//	Encrypt the private key using protection key
	var protectionKeyForSealing [32]byte
	copy(protectionKeyForSealing[:], protectionKeyBytes)
	encryptedPrivateKeyBytes := SealSymmetrically(privateKeyBytes[:], protectionKeyForSealing)

	//	Generate random 32 byte salt
	saltBytes, er := globalCommons.GenerateRandomBytes(commons.KEY_BYTES)
	if er != nil {
		return nil, errors.New(er, "Failed to generate salt for protection key", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	//	Safegaurd the protection key using Argon2i password based hashing.
	passwordDerivedKey := argon2.Key([]byte(password), saltBytes, 3, commons.KEY_BYTES*1024, 4, commons.KEY_BYTES)

	//	Encrypt the protection key using password derived key
	var passwordDerivedKeyForSealing [32]byte
	copy(passwordDerivedKeyForSealing[:], passwordDerivedKey)
	encryptedProtectionKeyBytes := SealSymmetrically(protectionKeyBytes, passwordDerivedKeyForSealing)

	return &commons.IssueKeyPairResponse{
		PublicKey:           publicKeyBytes[:],
		PrivateKey:          encryptedPrivateKeyBytes,
		DecryptedPrivateKey: privateKeyBytes[:],
		ProtectedKey:        encryptedProtectionKeyBytes,
		Salt:                saltBytes,
	}, nil
}

// Decrypt the org's symmetric key with your local public-private key.
func DecryptAsymmetricallyAnonymous(public, private, org_key []byte) ([]byte, *errors.Error) {

	var publicKey, privateKey [commons.KEY_BYTES]byte
	copy(publicKey[:], public)
	copy(privateKey[:], private)
	result, err := OpenAsymmetricallyAnonymous(org_key, publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetOrgKeyServerCopy(ctx context.ServiceContext, org_id string) ([]byte, *errors.Error) {

	errMessage := "Failed to get server-copy of org's encryption key"

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
	serverPrivateKey, er := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PRIVATE_KEY"))
	if er != nil {
		return nil, errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	serverPublicKey, er := base64.StdEncoding.DecodeString(os.Getenv("SERVER_PUBLIC_KEY"))
	if er != nil {
		return nil, errors.New(er, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	return DecryptAsymmetricallyAnonymous(serverPublicKey, serverPrivateKey, serverCopy)
}
