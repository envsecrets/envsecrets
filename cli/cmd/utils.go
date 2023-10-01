package cmd

import (
	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/internal/keys"
)

func Encrypt() {

	if commons.Secret.EnvID == "" {
		return
	}

	//	Decrypt the organisation key.
	var orgKey [32]byte
	decryptedOrgKey, err := keys.DecryptAsymmetricallyAnonymous(commons.KeysConfig.Public, commons.KeysConfig.Private, commons.ProjectConfig.Key)
	if err != nil {
		log.Debug(err)
		log.Fatal("Failed to decrypt the organisation's encryption key")
	}
	copy(orgKey[:], decryptedOrgKey)

	//	Encrypt the secrets
	if err := commons.Secret.Encrypt(orgKey); err != nil {
		log.Debug(err)
		log.Fatal("Failed to encrypt secrets")
	}
}

func Decrypt() {

	//	Decrypt the organisation key.
	var orgKey [32]byte
	decryptedOrgKey, err := keys.DecryptAsymmetricallyAnonymous(commons.KeysConfig.Public, commons.KeysConfig.Private, commons.ProjectConfig.Key)
	if err != nil {
		log.Debug(err)
		log.Fatal("Failed to decrypt the organisation's encryption key")
	}
	copy(orgKey[:], decryptedOrgKey)

	//	Encrypt the secrets
	if err := commons.Secret.Decrypt(orgKey); err != nil {
		log.Debug(err)
		log.Fatal("Failed to decrypt the secret")
	}
}

func DecryptAndDecode() {

	if commons.Secret.EnvID == "" {
		return
	}

	//	Decrypt the common secret.
	Decrypt()

	//	Decode the values.
	if err := commons.Secret.Decode(); err != nil {
		log.Debug(err)
		log.Fatal("Failed to decode the secret")
	}
}
