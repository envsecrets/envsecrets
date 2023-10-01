package internal

import (
	"fmt"
	"net/http"

	"github.com/envsecrets/envsecrets/cli/commons"
	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
)

func GetValues(ctx context.ServiceContext, client *clients.HTTPClient, options *GetValuesOptions) (*secretCommons.Secret, error) {

	req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodGet, commons.API+"/v1/secrets/values", nil)
	if err != nil {
		return nil, err
	}

	//	Initialize the query values.
	query := req.URL.Query()

	if options.Key != nil {
		query.Set("key", *options.Key)
	}

	if options.Version != nil {
		query.Set("version", fmt.Sprint(*options.Key))
	}

	//	If the environment token is passed,
	//	create a new HTTP client and attach it in the header.
	if options.Token != "" {

		client = clients.NewHTTPClient(&clients.HTTPConfig{
			BaseURL: commons.API + "/v1",
			CustomHeaders: []clients.CustomHeader{
				{
					Key:   string(clients.TokenHeader),
					Value: options.Token,
				},
			},
		})

	} else {

		//	Set the environment ID in query.
		query.Set("env_id", options.EnvID)
	}

	req.URL.RawQuery = query.Encode()

	var response clients.APIResponse
	if err := client.Run(commons.DefaultContext, req, &response); err != nil {
		return nil, err
	}

	var data secretCommons.Secret
	if err := globalCommons.MapToStruct(response.Data, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
