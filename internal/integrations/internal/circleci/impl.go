package circleci

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
		Type: clients.HTTPClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Circle-Token",
				Value: options.Credentials["token"].(string),
			},
		},
	})

	req, err := http.NewRequest(http.MethodGet, "https://circleci.com/api/v2/me/collaborations", nil)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func ListSubEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.HTTPClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Circle-Token",
				Value: options.Credentials["token"].(string),
			},
		},
	})

	req, err := http.NewRequest(http.MethodGet, "https://circleci.com/api/v2/pipeline"+"?org-slug="+options.OrgSlug, nil)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type: clients.HTTPClientType,
		CustomHeaders: []clients.CustomHeader{
			{
				Key:   "Circle-Token",
				Value: options.Credentials["token"].(string),
			},
		},
	})

	for key, payload := range *options.Data {
		body, err := json.Marshal(map[string]interface{}{
			"name":  key,
			"value": payload.Value,
		})
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://circleci.com/api/v2/project/%s/envvar", options.EntityDetails["project_slug"].(string)), bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		if err := client.Run(ctx, req, nil); err != nil {
			return err
		}
	}

	return nil
}
