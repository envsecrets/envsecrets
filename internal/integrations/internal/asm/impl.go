package asm

import (
	"encoding/base64"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/graphql"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func Setup(ctx context.ServiceContext, gqlClient *clients.GQLClient, options *SetupOptions) (*commons.Integration, error) {

	//	Encrypt the credentials
	credentials, err := commons.EncryptCredentials(ctx, options.OrgID, map[string]interface{}{
		"role_arn": options.RoleARN,
		"region":   options.Region,
	})
	if err != nil {
		return nil, err
	}

	//	Create a new record in Hasura.
	return graphql.Insert(ctx, gqlClient, &commons.AddIntegrationOptions{
		OrgID:       options.OrgID,
		Type:        commons.ASM,
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

func Sync(ctx context.ServiceContext, options *SyncOptions) (*secretsmanager.CreateSecretOutput, error) {

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

	//	Marshal the secrets
	payload, err := options.Secret.ToMap().Marshal()
	if err != nil {
		return nil, err
	}

	resp, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(options.EntityDetails["name"].(string)),
		SecretString: aws.String(string(payload)),
	})
	if err != nil {

		//	TODO: Use better error management over here.
		//	If the secret already exists, it returns HTTP Status Code 400.
		//	Use that response code for error handling.
		_, err = client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(options.EntityDetails["role_arn"].(string)),
			SecretString: aws.String(string(payload)),
		})
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
