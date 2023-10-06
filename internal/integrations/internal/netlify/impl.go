package netlify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
)

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
	for key, payload := range *options.Data {
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

	var response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	if err := client.Run(ctx, req, &response); err != nil {
		return err
	}

	if response.Code != http.StatusCreated {
		return fmt.Errorf("failed to sync the secrets to netlify site: %s", response.Message)
	}

	return nil
}
