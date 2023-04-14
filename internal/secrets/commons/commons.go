package commons

type Header string

const (
	VAULT_TOKEN Header = "X-Vault-Token"

	ENV_ID = "env_id"
)

const (
	VAULT_ROOT_TOKEN = "VAULT_ROOT_TOKEN"
)

const (
	Plaintext  Type = "plaintext"
	Ciphertext Type = "ciphertext"
)
