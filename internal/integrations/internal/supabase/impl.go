package supabase

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
)

func Setup(ctx context.ServiceContext, gqlClient *clients.GQLClient, options *SetupOptions) (*commons.Integration, *errors.Error) {

	//	Encrypt the credentials
	credentials, err := commons.EncryptCredentials(ctx, options.OrgID, map[string]interface{}{
		"token": options.Token,
	})
	if err != nil {
		return nil, err
	}

	//	Create a new record in Hasura.
	return graphql.Insert(ctx, gqlClient, &commons.AddIntegrationOptions{
		OrgID:       options.OrgID,
		Type:        commons.Supabase,
		Credentials: base64.StdEncoding.EncodeToString(credentials),
	})
}

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, *errors.Error) {

	errMessage := "Failed to fetch list of projects"

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: "Bearer " + options.Credentials["token"].(string),
	})

	req, er := http.NewRequest(http.MethodGet, "https://api.supabase.com/v1/projects", nil)
	if er != nil {
		return nil, errors.New(er, errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	var response []Project
	if err := client.Run(ctx, req, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) *errors.Error {

	errMessage := "Failed to sync secrets"

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: "Bearer " + options.Credentials["token"].(string),
	})

	body, er := json.Marshal(transform(options.Data))
	if er != nil {
		return errors.New(er, errMessage, errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	req, er := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.supabase.com/v1/projects/%s/secrets", options.EntityDetails["id"].(string)), bytes.NewBuffer(body))
	if er != nil {
		return errors.New(er, errMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	if err := client.Run(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
