package commons

type Header string

const (
	VAULT_TOKEN Header = "X-Vault-Token"
)

const (
	VAULT_ROOT_TOKEN = "VAULT_ROOT_TOKEN"
)

type Key string

const (
	ECDSA_P256   Key = "ecdsa-p256"
	AES256_GCM96 Key = "aes256-gcm96"
)
