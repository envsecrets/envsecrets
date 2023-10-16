package internal

import (
	"encoding/base64"
	"net/http"

	"github.com/envsecrets/envsecrets/cli/clients"
	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/internal/context"
)

func GetPublicKey(ctx context.ServiceContext, client *clients.HTTPClient, email string) ([]byte, error) {

	req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodGet, clients.API+"/v1/keys/public-key", nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return result, nil
}
