package commons

import (
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys"
)

func EncryptCredentials(ctx context.ServiceContext, org_id string, payload map[string]interface{}) ([]byte, *errors.Error) {

	errMessage := "Failed to encrypt credentials"

	//	Prepare the credentials
	credentials, er := json.Marshal(payload)
	if er != nil {
		return nil, errors.New(er, errMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Get the server's copy of organisation's encryption key.
	var orgKey [32]byte
	orgKeyBytes, err := keys.GetOrgKeyServerCopy(ctx, org_id)
	if err != nil {
		return nil, err
	}
	copy(orgKey[:], orgKeyBytes)

	//	Encrypt the secrets with the server-copy of organisation's encryption key.
	encrypted, err := keys.SealSymmetrically(credentials, orgKey)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func DecryptCredentials(ctx context.ServiceContext, org_id string, payload []byte) ([]byte, *errors.Error) {

	//	Get the server's copy of organisation's encryption key.
	var orgKey [32]byte
	orgKeyBytes, err := keys.GetOrgKeyServerCopy(ctx, org_id)
	if err != nil {
		return nil, err
	}
	copy(orgKey[:], orgKeyBytes)

	//	Decrypt the value using org-key.
	return keys.OpenSymmetrically(payload, orgKey)
}
