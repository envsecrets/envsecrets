package commons

type Type string

const (
	ECDSA_P256   Type = "ecdsa-p256"
	AES256_GCM96 Type = "aes256-gcm96"
)

type KeyType string

const (
	EncryptionKey KeyType = "encryption-key"
	SigningKey    KeyType = "signing-key"
	HmacKey       KeyType = "hmac-key"
)
