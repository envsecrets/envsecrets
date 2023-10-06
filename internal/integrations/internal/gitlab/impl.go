package gitlab

import (
	"fmt"
	"os"

	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
)

// Prepares credentials to be saved in the database.
func PrepareCredentials(ctx context.ServiceContext, options *PrepareCredentialsOptions) (map[string]interface{}, error) {

	//	Exchange the code for Access Token
	response, err := GetAccessToken(ctx, &TokenRequestOptions{
		Code:        options.Code,
		RedirectURI: os.Getenv("REDIRECT_DOMAIN") + "/v1/integrations/gitlab/callback/setup",
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token_type":    response.TokenType,
		"refresh_token": response.RefreshToken,
	}, nil
}

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	//	Refresh access token
	access, err := RefreshToken(ctx, &TokenRefreshOptions{
		RefreshToken:  options.Credentials["refresh_token"].(string),
		OrgID:         options.OrgID,
		IntegrationID: options.IntegrationID,
	})
	if err != nil {
		return nil, err
	}

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: fmt.Sprintf("%s %s", access.TokenType, access.AccessToken),
	})

	switch options.Type {
	case ProjectType:
		return ListProjects(ctx, client)
	case GroupType:
		return ListGroups(ctx, client)
	}

	return nil, errors.New("invalid entity type")
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	Refresh access token
	access, err := RefreshToken(ctx, &TokenRefreshOptions{
		RefreshToken:  options.Credentials["refresh_token"].(string),
		OrgID:         options.OrgID,
		IntegrationID: options.IntegrationID,
	})
	if err != nil {
		return err
	}

	//	Initialize a new HTTP client.
	client := clients.NewHTTPClient(&clients.HTTPConfig{
		Type:          clients.HTTPClientType,
		Authorization: fmt.Sprintf("%s %s", access.TokenType, access.AccessToken),
	})

	for key, payload := range *options.Data {

		switch EntityType(options.EntityDetails["type"].(string)) {
		case ProjectType:
			if _, err := CreateProjectVariable(ctx, client, &CreateVariableOptions{
				ID: int64(options.EntityDetails["id"].(float64)),
				Variable: Variable{
					Key:   key,
					Value: fmt.Sprint(payload.Value),
				},
			}); err != nil {
				return err
			}
		case GroupType:
			if _, err := CreateGroupVariable(ctx, client, &CreateVariableOptions{
				ID: int64(options.EntityDetails["id"].(float64)),
				Variable: Variable{
					Key:   key,
					Value: fmt.Sprint(payload.Value),
				},
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
