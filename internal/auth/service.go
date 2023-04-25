package auth

import (
	"bytes"
	"encoding/base64"
	internalErrors "errors"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/auth/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/nhost"
	"github.com/envsecrets/envsecrets/internal/users"
)

func Signup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SignupOptions) *errors.Error {

	//	Signup on Nhost
	if err := nhost.Signup(ctx, &nhost.SignupOptions{
		Email:    options.Email,
		Password: options.Password,
		Options: map[string]interface{}{
			"displayName": options.Name,
		},
	}); err != nil {
		return errors.New(internalErrors.New(err.Message), "Failed to signup the user", errors.ErrorTypeBadResponse, errors.ErrorSourceNhost)
	}

	//	Fetch the user with their email.
	user, err := users.GetByEmail(ctx, client, options.Email)
	if err != nil {
		return err
	}

	//	Generate Key pair
	pair, err := keys.GenerateKeyPair(options.Password)
	if err != nil {
		return err
	}

	//	Upload the keys to their cloud account.
	if err := keys.CreateWithUserID(ctx, client, &keyCommons.CreateWithUserIDOptions{
		PublicKey:    base64.StdEncoding.EncodeToString(pair.PublicKey),
		PrivateKey:   base64.StdEncoding.EncodeToString(pair.PrivateKey),
		ProtectedKey: base64.StdEncoding.EncodeToString(pair.ProtectedKey),
		Salt:         base64.StdEncoding.EncodeToString(pair.Salt),
		UserID:       user.ID,
	}); err != nil {
		return err
	}

	return nil
}

func UpdatePassword(ctx context.ServiceContext, client *clients.HTTPClient, options *commons.UpdatePasswordOptions) *errors.Error {

	errMessage := "Failed to update password"
	body, err := options.Marshal()
	if err != nil {
		return errors.New(err, errMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodPost, os.Getenv("NHOST_AUTH_URL")+"/user/password", bytes.NewBuffer(body))
	if err != nil {
		return errors.New(err, errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return client.Run(ctx, req, nil)
}
