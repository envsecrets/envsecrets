package netlify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
)

func Setup(ctx context.ServiceContext, gqlClient *clients.GQLClient, options *SetupOptions) (*commons.Integration, error) {

	//	Encrypt the credentials
	credentials, err := commons.EncryptCredentials(ctx, options.OrgID, options.toMap())
	if err != nil {
		return nil, err
	}

	//	Create a new record in Hasura.
	return graphql.Insert(ctx, gqlClient, &commons.AddIntegrationOptions{
		OrgID:       options.OrgID,
		Type:        commons.Netlify,
		Credentials: base64.StdEncoding.EncodeToString(credentials),
	})
}

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: "Bearer " + options.Credentials["token"].(string),
	})

	req, err := http.NewRequest(http.MethodGet, "https://api.netlify.com/api/v1/sites", nil)
	if err != nil {
		return nil, err
	}

	var response []Site
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: "Bearer " + options.Credentials["token"].(string),
	})

	//	Fetch the account ID from netlify
	user, err := fetchAccounts(ctx, client)
	if err != nil {
		return err
	}

	var result []map[string]interface{}
	for key, payload := range options.Secret.Data {
		result = append(result, map[string]interface{}{
			"key":    key,
			"scopes": []string{"builds", "functions", "runtime", "post-processing"},
			"values": []map[string]interface{}{
				{
					"value":   payload.Value,
					"context": "all",
				},
			},
		})
	}

	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	/* 	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.netlify.com/api/v1/accounts/%s/env?site_id=%s", user.ID, options.EntityDetails["id"].(string)), bytes.NewBuffer(body))
	   	if err != nil {
	   		return errors.New(er, errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	   	}
	*/

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.netlify.com/api/v1/accounts/%s/env", user.ID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	var response interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return err
	}

	return nil
}
