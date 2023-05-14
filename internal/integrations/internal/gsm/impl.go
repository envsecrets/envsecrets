package gsm

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"
	"google.golang.org/api/option"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func Setup(ctx context.ServiceContext, gqlClient *clients.GQLClient, options *SetupOptions) (*commons.Integration, *errors.Error) {

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

func ListEntities(ctx context.ServiceContext, options *ListOptions) (interface{}, *errors.Error) {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(options.Credentials["region"].(string)),
	)
	if err != nil {
		return nil, errors.New(err, "Failed to push secrets to ASM", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
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
		return nil, errors.New(err, "Failed to list secrets from ASM", errors.ErrorTypeBadRequest, errors.ErrorSourceHTTP)
	}

	return resp, nil
}

func Sync(ctx context.ServiceContext, options *SyncOptions) *errors.Error {

	//	TODO: Inform the user in case of failed synchronization
	//	Take email from GSM service account credentials.
	//	Ex: "client_email": "doppler-secret-manager@envsecrets.iam.gserviceaccount.com"

	//	Marshal the credentials
	creds, er := json.Marshal(options.Credentials)
	if er != nil {
		return errors.New(er, "Failed to sync secrets to GSM", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	// Create the client.
	client, er := secretmanager.NewClient(ctx, option.WithCredentialsJSON(creds))
	if er != nil {
		// The most likely causes of the error are:
		//     1 - google application creds failed
		//     2 - secret already exists
		return errors.New(er, "Failed to sync secrets to GSM", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}
	defer client.Close()

	PARENT := fmt.Sprintf("projects/%v", options.Credentials["project_id"])

	//	Prepare the payload
	payload, er := json.Marshal(toKVMap(options.Data))
	if er != nil {
		return errors.New(er, "Failed to sync secrets to GSM", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Create a new version
	_, er = AddSecretVersion(ctx, client, fmt.Sprintf("%s/secrets/%s", PARENT, options.EntityDetails["name"].(string)), payload)
	if er != nil {

		//	Create the secret if it doesn't exist
		if strings.Contains(er.Error(), "NotFound") {

			if _, er = CreateSecret(ctx, client, PARENT, options.EntityDetails["name"].(string)); er != nil {
				return errors.New(er, "Failed to sync secrets to GSM", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
			}

			return Sync(ctx, options)
		}

		return errors.New(er, "Failed to sync secrets to GSM", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return nil
}
