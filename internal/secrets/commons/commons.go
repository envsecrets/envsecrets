package commons

type Header string

const (
	VAULT_TOKEN Header = "X-Vault-Token"
)

const (
	VAULT_ROOT_TOKEN = "VAULT_ROOT_TOKEN"
)

const (
	Plaintext  Type = "plaintext"
	Ciphertext Type = "ciphertext"
)
