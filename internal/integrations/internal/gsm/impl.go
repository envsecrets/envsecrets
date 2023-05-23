package gsm

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	"google.golang.org/api/option"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func Setup(ctx context.ServiceContext, gqlClient *clients.GQLClient, options *SetupOptions) (*commons.Integration, error) {

	//	Encrypt the credentials
	credentials, err := commons.EncryptCredentials(ctx, options.OrgID, options.Keys)
	if err != nil {
		return nil, err
	}

	//	Create a new record in Hasura.
	return graphql.Insert(ctx, gqlClient, &commons.AddIntegrationOptions{
		OrgID:       options.OrgID,
		Type:        commons.GSM,
		Credentials: base64.StdEncoding.EncodeToString(credentials),
	})
}

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, error) {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(options.Credentials["region"].(string)),
	)
	if err != nil {
		return nil, err
	}

	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, options.Credentials["role_arn"].(string), func(aro *stscreds.AssumeRoleOptions) {
		aro.RoleARN = options.Credentials["role_arn"].(string)
		aro.ExternalID = aws.String(options.OrgID)
	})
	cfgCopy := cfg.Copy()
	cfgCopy.Credentials = aws.NewCredentialsCache(provider)
	client := secretsmanager.NewFromConfig(cfg)

	resp, err := client.ListSecrets(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) error {

	//	TODO: Inform the user in case of failed synchronization
	//	Take email from GSM service account credentials.
	//	Ex: "client_email": "doppler-secret-manager@envsecrets.iam.gserviceaccount.com"

	//	Marshal the credentials
	creds, err := json.Marshal(options.Credentials)
	if err != nil {
		return err
	}

	// Create the client.
	client, err := secretmanager.NewClient(ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		// The most likely causes of the error are:
		//     1 - google application creds failed
		//     2 - secret already exists
		return err
	}
	defer client.Close()

	PARENT := fmt.Sprintf("projects/%v", options.Credentials["project_id"])

	//	Prepare the payload
	payload, err := options.Secrets.ToMap().Marshal()
	if err != nil {
		return err
	}

	//	Create a new version
	_, err = AddSecretVersion(ctx, client, fmt.Sprintf("%s/secrets/%s", PARENT, options.EntityDetails["name"].(string)), payload)
	if err != nil {

		//	Create the secret if it doesn't exist
		if strings.Contains(err.Error(), "NotFound") {

			if _, err = CreateSecret(ctx, client, PARENT, options.EntityDetails["name"].(string)); err != nil {
				return err
			}

			return Sync(ctx, options)
		}

		return err
	}

	return nil
}
