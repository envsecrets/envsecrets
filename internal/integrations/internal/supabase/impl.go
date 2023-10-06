package supabase

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

	req, err := http.NewRequest(http.MethodGet, "https://api.supabase.com/v1/projects", nil)
	if err != nil {
		return nil, err
	}

	var response []Project
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

	var result []map[string]interface{}
	for key, payload := range *options.Data {
		result = append(result, map[string]interface{}{
			"name":  key,
			"value": payload.Value,
		})
	}

	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.supabase.com/v1/projects/%s/secrets", options.EntityDetails["id"].(string)), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if err := client.Run(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
