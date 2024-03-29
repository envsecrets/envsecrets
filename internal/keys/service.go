package keys

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/keys/graphql"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/utils"
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

func UpdateSyncKey(ctx context.ServiceContext, client *clients.GQLClient, key_id string) ([]byte, error) {

	//	Generate a separate random symmetric key
	syncKeyBytes, err := utils.GenerateRandomBytes(commons.KEY_BYTES)
	if err != nil {
		return nil, err
	}

	//	Encrypt the sync key using server's symmetric key.
	encryptedSyncKeyBytes, err := SealSymmetricallyByServer(syncKeyBytes)
	if err != nil {
		return nil, err
	}

	if err := graphql.UpdateSyncKey(ctx, client, &commons.UpdateSyncKeyOptions{
		KeyID:   key_id,
		SyncKey: base64.StdEncoding.EncodeToString(encryptedSyncKeyBytes),
	}); err != nil {
		return nil, err
	}

	return encryptedSyncKeyBytes, nil
}

func GetByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) (*commons.Key, error) {
	return graphql.GetByUserID(ctx, client, user_id)
}

func GetPublicKey(ctx context.ServiceContext, client *clients.GQLClient) ([]byte, error) {
	return graphql.GetPublicKey(ctx, client)
}

func GetPublicKeyByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) ([]byte, error) {
	return graphql.GetPublicKeyByUserID(ctx, client, user_id)
}

func GetSyncKey(ctx context.ServiceContext, client *clients.GQLClient) ([]byte, error) {
	return graphql.GetSyncKey(ctx, client)
}

func GetSyncKeyByUserID(ctx context.ServiceContext, client *clients.GQLClient, user_id string) ([]byte, error) {
	return graphql.GetSyncKeyByUserID(ctx, client, user_id)
}

func GetPublicKeyByUserEmail(ctx context.ServiceContext, client *clients.GQLClient, email string) ([]byte, error) {
	return graphql.GetPublicKeyByUserEmail(ctx, client, email)
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

func SealSymmetricallyByServer(message []byte) ([]byte, error) {

	//	Encrypt the sync key using server's symmetric key
	encodedServerKey := os.Getenv("SERVER_SYMMETRIC_KEY")
	if encodedServerKey == "" {
		return nil, errors.New("SERVER_SYMMETRIC_KEY is not set")
	}

	decodedServerKey, err := base64.StdEncoding.DecodeString(encodedServerKey)
	if err != nil {
		return nil, err
	}

	var serverKey [32]byte
	copy(serverKey[:], decodedServerKey)

	return SealSymmetrically(message[:], serverKey)
}

func OpenSymmetricallyByServer(payload []byte) ([]byte, error) {

	//	Encrypt the sync key using server's symmetric key
	encodedServerKey := os.Getenv("SERVER_SYMMETRIC_KEY")
	if encodedServerKey == "" {
		return nil, commons.ErrNoServerKey
	}

	decodedServerKey, err := base64.StdEncoding.DecodeString(encodedServerKey)
	if err != nil {
		return nil, err
	}

	var serverKey [32]byte
	copy(serverKey[:], decodedServerKey)

	return OpenSymmetrically(payload[:], serverKey)
}

func GenerateKeyPair(password string) (*commons.IssueKeyPairResponse, error) {

	publicKeyBytes, privateKeyBytes, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	//	Generate a separate random symmetric key
	syncKeyBytes, err := utils.GenerateRandomBytes(commons.KEY_BYTES)
	if err != nil {
		return nil, err
	}

	//	Generate a separate random symmetric key
	protectionKeyBytes, err := utils.GenerateRandomBytes(commons.KEY_BYTES)
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
	saltBytes, err := utils.GenerateRandomBytes(commons.KEY_BYTES)
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
		SyncKey:             syncKeyBytes,
	}, nil
}

func DecryptPayload(payload *commons.Payload, password string) error {

	//	Regenerate the key from user's password
	passwordDerivedKey := argon2.Key([]byte(password), payload.Salt, 3, commons.KEY_BYTES*1024, 4, commons.KEY_BYTES)

	//	Use the password derived key to decrypt protection key.
	var passwordDerivedKeyForOpening [32]byte
	copy(passwordDerivedKeyForOpening[:], passwordDerivedKey)
	protectionKey, err := OpenSymmetrically(payload.ProtectedKey, passwordDerivedKeyForOpening)
	if err != nil {
		return err
	}

	payload.ProtectedKey = protectionKey

	//	Decrypt the private key using the decrypted protection key.
	var protectionKeyForOpening [32]byte
	copy(protectionKeyForOpening[:], protectionKey)
	privateKey, err := OpenSymmetrically(payload.PrivateKey, protectionKeyForOpening)
	if err != nil {
		return err
	}

	payload.PrivateKey = privateKey

	//	Decrypt the sync key using the server's own encryption key.
	syncKey, err := OpenSymmetricallyByServer(payload.SyncKey)
	if err != nil && err != commons.ErrNoServerKey {
		return err
	}

	payload.SyncKey = syncKey

	return nil
}

func DecryptMemberKey(ctx context.ServiceContext, client *clients.GQLClient, user_id string, options *commons.DecryptOptions) ([]byte, error) {

	//	Get the user's key pair.
	keyPair, err := GetByUserID(ctx, client, user_id)
	if err != nil {
		return nil, err
	}

	//	Base64 decode the key pair.
	payload, err := keyPair.Decode()
	if err != nil {
		return nil, err
	}

	//	Decrypt the user's private key.
	if err := DecryptPayload(payload, options.Password); err != nil {
		return nil, err
	}

	//	Pull the user's copy of organisation's encryption key.
	userKey, err := memberships.GetKey(ctx, client, &memberships.GetKeyOptions{
		OrgID:  options.OrgID,
		UserID: user_id,
	})
	if err != nil {
		return nil, err
	}

	decryptedKey, err := DecryptAsymmetricallyAnonymous(payload.PublicKey, payload.PrivateKey, userKey)
	if err != nil {
		return nil, err
	}

	return decryptedKey, nil
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
