package internal

import (
	"encoding/base64"
	"net/http"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
)

func GetPublicKey(ctx context.ServiceContext, client *clients.HTTPClient, email string) ([]byte, *errors.Error) {

	errMessage := "Failed to get public key"

	req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodGet, commons.API+"/v1/keys/public-key", nil)
	if err != nil {
		return nil, errors.New(err, errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	//	Initialize the query values.
	query := req.URL.Query()
	query.Set("email", email)

	req.URL.RawQuery = query.Encode()

	var response clients.APIResponse
	if err := client.Run(commons.DefaultContext, req, &response); err != nil {
		return nil, err
	}

	result, err := base64.StdEncoding.DecodeString(response.Data.(string))
	if err != nil {
		return nil, errors.New(err, errMessage, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
	}

	return result, nil
}
